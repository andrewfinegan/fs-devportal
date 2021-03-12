package doc

import (
	"bytes"
	"context"
	"devportal/api/product"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/v33/github"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

//Get Docs from Github by passing the document path in request body
func GetDocumentServiceSelector(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	docPath := query.Get("path")

	fmt.Println("Doc Path ::", docPath)

	if docPath == "" {
		GetDocsTree(w, r)
	} else {
		GetDocument(w, r, docPath)
	}
}

//Get Docs from Github by passing the document path in request body
func GetDocument(w http.ResponseWriter, r *http.Request, docPath string) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	//Get Docs By Path
	response, err := http.Get(product.DevPortalConfig.GitHub.GitHubContentFullPath + docPath)

	if err != nil {
		w.Write([]byte("Markdown file not found"))
		w.WriteHeader(http.StatusBadRequest)
	} else {
		buf := new(bytes.Buffer)
		buf.ReadFrom(response.Body)
		responseByte := buf.Bytes()
		defer response.Body.Close()

		w.Write(responseByte)
		w.WriteHeader(http.StatusOK)
	}

}

//Get DocTree from Github for /docs folder
func GetDocsTree(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if product.DevPortalConfig.GitHub.GitHubSourceOwner == "" || product.DevPortalConfig.GitHub.GitHubSourceRepo == "" || product.DevPortalConfig.GitHub.GitHubContentBranch == "" {
		log.Fatal("You need to specify a non-empty value for `-source-owner`, `-source-repo`, and `-DevPortalConfig.GitHub.GitHubBaseBranch`")
	}

	//Basic Auth
	tp := github.BasicAuthTransport{
		Username: product.DevPortalConfig.GitHub.GitHubUserName,
		Password: product.DevPortalConfig.GitHub.GitHubAuthToken,
	}

	client = github.NewClient(tp.Client())

	ref, err := getRef()
	if err != nil {
		log.Printf("Unable to get/create the commit reference: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ref == nil {
		log.Println("No error where returned but the reference is nil")
		return
	}

	tree, err := getTree(ref)
	if err != nil {
		log.Printf("Unable to create the tree based on the provided files: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if tree == nil {
		log.Println("No error where returned but doc tree is nil")
		return
	}

	var respBytes []byte
	respBytes, err = json.Marshal(tree)
	if err != nil {
		log.Printf("Unable marshall tree structure: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = w.Write(respBytes)
	if err != nil {
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)

}

var client *github.Client
var ctx = context.Background()

// getRef returns the commit branch reference object if it exists or creates it
// from the base branch before returning it.
func getRef() (ref *github.Reference, err error) {
	if ref, _, err = client.Git.GetRef(ctx, product.DevPortalConfig.GitHub.GitHubSourceOwner, product.DevPortalConfig.GitHub.GitHubSourceRepo, "refs/heads/"+product.DevPortalConfig.GitHub.GitHubContentBranch); err == nil {
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
	tree, _, err := client.Git.GetTree(ctx, product.DevPortalConfig.GitHub.GitHubSourceOwner, product.DevPortalConfig.GitHub.GitHubSourceRepo, *ref.Object.SHA, true)
	var newDoc Document
	docTree = new(DocumentTree)

	for _, fileArg := range tree.Entries {
		if strings.Contains(*fileArg.Path, "docs/") && isMDFileOrFolder(fileArg) {
			newDoc, err = buildDoc(fileArg)
			docTree.Documents = append(docTree.Documents, newDoc)
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
