// ################################################################################
// Copyright © 2021-2022 Fiserv, Inc. or its affiliates. 
// Fiserv is a trademark of Fiserv, Inc., 
// registered or used in the United States and foreign countries, 
// and may or may not be registered in your country.  
// All trademarks, service marks, 
// and trade names referenced in this 
// material are the property of their 
// respective owners. This work, including its contents 
// and programming, is confidential and its use 
// is strictly limited. This work is furnished only 
// for use by duly authorized licensees of Fiserv, Inc. 
// or its affiliates, and their designated agents 
// or employees responsible for installation or 
// operation of the products. Any other use, 
// duplication, or dissemination without the 
// prior written consent of Fiserv, Inc. 
// or its affiliates is strictly prohibited. 
// Except as specified by the agreement under 
// which the materials are furnished, Fiserv, Inc. 
// and its affiliates do not accept any liabilities 
// with respect to the information contained herein 
// and are not responsible for any direct, indirect, 
// special, consequential or exemplary damages 
// resulting from the use of this information. 
// No warranties, either express or implied, 
// are granted or extended by this work or 
// the delivery of this work
// ################################################################################

package doc

import (
	"bytes"
	"context"
	"devportal/api"
	"devportal/config"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v33/github"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()

// Adding function to check if on given dir any docignore file exists. 
func checkForDocIgnoreFile( baseUrl string)([]string){
	var txtlines []string
	docIgnoreURL := baseUrl + "/.docignore" 
	response, err := http.Get(docIgnoreURL) 
 
	if (err != nil){
		config.Logger.Error("docignore file Not found at requested dir")
	} 
	if response.StatusCode == http.StatusOK {

	/* Note: io.ReadAll is only compatible wiht Go 1.16 or later. */
	
			bodyContent, err := ioutil.ReadAll(response.Body)
			if err != nil {
				config.Logger.Error("Not able to read data from docignore file")
			} 
	scanner := bufio.NewScanner(strings.NewReader(string(bodyContent)))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
	// skipping any comments in docignore file	
	statementCheck := strings.HasPrefix(scanner.Text(), "#")
		if (!statementCheck){
			if (len(scanner.Text()) > 0 ){
				txtlines = append(txtlines, scanner.Text())
			} 
		} 
	}
	}
	return txtlines
}

//Get Docs from Github by passing the document path in request body
func GetDocumentServiceSelector(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	docPath := query.Get("path")
	documentFileUrl := config.AppConfig.GitHub.GitHubContentFullPath

	// Checking for docignore in a given url
	content := checkForDocIgnoreFile(documentFileUrl)
	if (content != nil){ 
			if docPath == "" {
				GetDocsTree(w, r)
			} else {
		 			var docExists bool
					for _, eachline := range content {
					// Checking directory(s) and if any exist, just removing asterisk and check if dir exist in requested doc path	
					docExists = strings.Contains(docPath ,strings.Replace(eachline , "**" , "",-1))
					if (docExists){
						break;
					}
				}
				if (docExists){
					api.WriteErrorResponse(w, api.ErrorResponse{
					StatusCode:   http.StatusNotFound,
					Message:      "Markdown file Ignored: " + docPath,
					ResponseCode: "Ignored"}, "")
				}else{
					// Spliting doc path file to its current directory to check if any docignore file exit at this level
					positionIndex  := strings.LastIndex(docPath, "/")  
					if positionIndex > -1 { 
						docPathDir := docPath[:positionIndex]
						documentDirPath := config.AppConfig.GitHub.GitHubContentFullPath + docPathDir + "/.docignore" 
						response, err := http.Get(documentDirPath) 
						if (err != nil){
							// Any error caused by retirieving docignore file, process the data doc path
							GetDocument(w, r, docPath)
						}
						// if docIgnore exist at this sub dir then will ignore the request
						if response.StatusCode == http.StatusOK {
							api.WriteErrorResponse(w, api.ErrorResponse{
							StatusCode:   http.StatusNotFound,
							Message:      "Markdown file Ignored: " + docPath,
							ResponseCode: "Ignored"}, "")
						}else{ 
							GetDocument(w, r, docPath)
						}  
					} 
				} 
	 		}
		}else{
			if docPath == "" {
				GetDocsTree(w, r)
			}else{
				GetDocument(w, r, docPath)
		}
	} 
}

//Get Docs from Github by passing the document path in request body
func GetDocument(w http.ResponseWriter, r *http.Request, docPath string) {
	//Get Docs By Path
	documentFileUrl := config.AppConfig.GitHub.GitHubContentFullPath + docPath
	logger.Println("Document File Path:", documentFileUrl)
	response, err := http.Get(documentFileUrl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if response.StatusCode != http.StatusOK {
		api.WriteErrorResponse(w, api.ErrorResponse{
			StatusCode:   http.StatusNotFound,
			Message:      "Markdown file not found: " + docPath,
			ResponseCode: "NotFound"}, "")
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(response.Body)
		responseByte := buf.Bytes()
		defer response.Body.Close()

		_, err = w.Write(responseByte)
		w.WriteHeader(http.StatusOK)
	}

}

var client *github.Client
var ctx = context.Background()

// getRef returns the commit branch reference object if it exists or creates it
// from the base branch before returning it.
func getRef() (ref *github.Reference, err error) {
	if ref, _, err = client.Git.GetRef(ctx, config.AppConfig.GitHub.GitHubSourceOwner, config.AppConfig.GitHub.GitHubSourceRepo, "refs/heads/"+config.AppConfig.GitHub.GitHubContentBranch); err == nil {
		return ref, nil
	}

	return ref, err
}

func buildDoc(treeEntry *github.TreeEntry) (doc Document, err error) {

	doc.Path = treeEntry.Path

	basename := filepath.Base(*doc.Path)
	name := strings.TrimSuffix(basename, filepath.Ext(basename))

	docTitle := name
	docTitle = strings.ReplaceAll(docTitle, "-", " ")
	doc.Title = &docTitle
	doc.Type = treeEntry.Type
	if *treeEntry.Type == "tree" {
		*doc.Type = "folder"
	} else {
		*doc.Type = "document"
	}
	return doc, err
}

func isMDFileOrFolder(treeEntry *github.TreeEntry) (result bool) {

	if *treeEntry.Type == "tree" {
		result = true
	} else {
		basename := filepath.Base(*treeEntry.Path)
		var ext = filepath.Ext(basename)
		if ext == ".md" {
			result = true
		}
	}
	return result

}

// getTree generates the document tree
// getTree generates the document tree
func getTree(ref *github.Reference) (docTree *DocumentTree, err error) {
	tree, _, err := client.Git.GetTree(ctx, config.AppConfig.GitHub.GitHubSourceOwner, config.AppConfig.GitHub.GitHubSourceRepo, *ref.Object.SHA, true)
	var newDoc Document
	docTree = new(DocumentTree)
	if tree != nil {

		documentFileUrl := config.AppConfig.GitHub.GitHubContentFullPath 
		// Checking for docignore in a given url
		content := checkForDocIgnoreFile(documentFileUrl) 
		for _, fileArg := range tree.Entries {
			if strings.Contains(*fileArg.Path, "docs/") && isMDFileOrFolder(fileArg) {
				newDoc, err = buildDoc(fileArg)
				if (content != nil){
					var docExists1 bool
					var docExists2 bool
					var docExists3 bool
					for _, eachline := range content {
					// Checking directory(s) and if any exist, just removing asterisk and check if dir exist in requested doc path	
					docExists1 = strings.Contains(*fileArg.Path ,strings.Replace(eachline , "**" , "",-1))
					docExists2 = strings.Contains(*fileArg.Path ,strings.Replace(eachline , "/**" , "",-1))
					docExists3 = strings.Contains(*fileArg.Path ,strings.Replace(eachline , "/" , "",-1))
					if (docExists1 || docExists2 || docExists3){
						break ;
					}
				}
				if (!docExists1 && !docExists2 && !docExists3){
					// Checking docignore file to exclude directory or file included into docignore file.
					docTree.Documents = append(docTree.Documents, newDoc)
				}
				} else{
					// If docignore file doesn't exit not filtering of any file(s) and dir(s).
					docTree.Documents = append(docTree.Documents, newDoc)
			 	}
			}
		}
	}
	return docTree, err
}

type Document struct {
	Path      *string    `json:"path,omitempty"`
	Type      *string    `json:"type,omitempty"`
	Title     *string    `json:"title,omitempty"`
	Documents []Document `json:"docs,omitempty"`
}

type DocumentTree struct {
	Documents []Document `json:"docs,omitempty"`
}

//Get DocTree from Github for /docs folder
func GetDocsTree(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if config.AppConfig.GitHub.GitHubSourceOwner == "" || config.AppConfig.GitHub.GitHubSourceRepo == "" || config.AppConfig.GitHub.GitHubContentBranch == "" {
		logger.Fatal("You need to specify a non-empty value for `-source-owner`, `-source-repo`, and `-AppConfig.GitHub.GitHubContentBranch`")
	}

	//Basic Auth
	tp := github.BasicAuthTransport{
		Username: config.AppConfig.GitHub.GitHubUserName,
		Password: config.AppConfig.GitHub.GitHubAuthToken,
	}

	client = github.NewClient(tp.Client())

	ref, err := getRef()
	if err != nil {
		api.WriteErrorResponse(w, api.ErrorResponse{
			StatusCode:   http.StatusInternalServerError,
			Message:      "Unable to build the tree structure",
			ResponseCode: "ServiceError"}, fmt.Sprintf("Unable to get/create the commit reference: %s\n", err))
		return
	}
	if ref == nil {
		api.WriteErrorResponse(w, api.ErrorResponse{
			StatusCode:   http.StatusInternalServerError,
			Message:      "Unable to build the tree structure",
			ResponseCode: "ServiceError"}, fmt.Sprintf("No error where returned but the reference is nil"))
		return
	}

	tree, err := getTree(ref)
	if err != nil {
		api.WriteErrorResponse(w, api.ErrorResponse{
			StatusCode:   http.StatusInternalServerError,
			Message:      "Unable to build the tree structure",
			ResponseCode: "ServiceError"}, fmt.Sprintf("Unable to create the tree based on the provided files: %s\n", err))
		return
	}

	if tree == nil {
		api.WriteErrorResponse(w, api.ErrorResponse{
			StatusCode:   http.StatusInternalServerError,
			Message:      "Unable to build the tree structure",
			ResponseCode: "ServiceError"}, fmt.Sprintf("No error where returned but doc tree is nil"))
		return
	}

	var respBytes []byte
	respBytes, err = json.Marshal(tree)
	if err != nil {
		api.WriteErrorResponse(w, api.ErrorResponse{
			StatusCode:   http.StatusInternalServerError,
			Message:      "Unable to build the tree structure",
			ResponseCode: "ServiceError"}, fmt.Sprintf("Unable marshall tree structure: %s\n", err))
		return
	}

	_, err = w.Write(respBytes)
	if err != nil {
		logger.Println(err)

		return
	}
	w.WriteHeader(http.StatusOK)

}
