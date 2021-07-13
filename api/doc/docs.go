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

//Get Docs from Github by passing the document path in request body
func GetDocumentServiceSelector(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	docPath := query.Get("path")

	if docPath == "" {
		GetDocsTree(w, r)
	} else {

		GetDocument(w, r, docPath)
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
func getTree(ref *github.Reference) (docTree *DocumentTree, err error) {
	tree, _, err := client.Git.GetTree(ctx, config.AppConfig.GitHub.GitHubSourceOwner, config.AppConfig.GitHub.GitHubSourceRepo, *ref.Object.SHA, true)
	var newDoc Document
	docTree = new(DocumentTree)

	if tree != nil {
		for _, fileArg := range tree.Entries {
			if strings.Contains(*fileArg.Path, "docs/") && isMDFileOrFolder(fileArg) {
				newDoc, err = buildDoc(fileArg)
				docTree.Documents = append(docTree.Documents, newDoc)
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
