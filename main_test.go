package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetBook(t *testing.T) {
	testCases := []struct {
		desc        string
		reqEndpoint string
		expRes      []Book
		statusCode  int
	}{
		{"success: get all books", "book", []Book{
			{1, "English", Author{Id: 1}, "scholostic", "11/9/2000"},
			{2, "Mathematics", Author{Id: 1}, "penguin", "9/11/1999"},
			{3, "Science", Author{Id: 1}, "Arihant", "3/4/1980"}},
			http.StatusOK},
		{"failure: wrong endpoint", "books", []Book{{}}, http.StatusBadRequest},
	}
	for j, tc := range testCases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "localhost:8000/"+tc.reqEndpoint, nil)
		getBooks(w, req)
		res, _ := io.ReadAll(w.Result().Body)
		resBooks := []Book{}
		err := json.Unmarshal(res, &resBooks)
		log.Println(err)
		if len(resBooks) != len(tc.expRes) {
			t.Errorf("%v test failed %v", j, tc.desc)
		}

		if reflect.DeepEqual(resBooks, tc.expRes) {
			t.Errorf("%v test failed %v", j, tc.desc)
		}

	}
}

func TestGetBookById(t *testing.T) {
	testCases := []struct {
		desc        string
		reqEndpoint string
		expRes      Book
		statusCode  int
	}{
		{"successfully retrieved", "1", Book{1, "English", Author{}, "scholostic", "13/9/2000"}, http.StatusOK},
		{"invalid id", "-5", Book{}, http.StatusBadRequest},
		{"id doesn't exist", "-1", Book{}, http.StatusBadRequest},
	}
	for i, tc := range testCases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "localhost:8000/book/"+tc.reqEndpoint, nil)
		getBookById(w, req)
		res, _ := io.ReadAll(w.Result().Body)
		resBook := Book{}
		json.Unmarshal(res, &resBook)

		if reflect.DeepEqual(resBook, tc.expRes) {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestPostBook(t *testing.T) {
	testCases := []struct {
		desc        string
		reqEndpoint string
		reqBody     Book
		expRes      Book
		statusCode  int
	}{
		{"successfully posted", "book", Book{Title: "the wall", Author: Author{Id: 1}, Publication: "Arihant", PublishedDate: "6/7/2017"}, Book{Title: "the wall", Author: Author{Id: 1}, Publication: "Arihant", PublishedDate: "6/7/2017"}, http.StatusCreated},
		{"successfully posted", "book", Book{Title: "city", Author: Author{Id: 1}, Publication: "Arihant", PublishedDate: "1/1/2001"}, Book{Title: "city", Author: Author{Id: 1}, Publication: "Arihant", PublishedDate: "1/1/2001"}, http.StatusCreated},
		{"failure : Invalid Publication", "book", Book{Title: "the wall", Author: Author{Id: 1}, Publication: "the sun", PublishedDate: "9/92012"}, Book{}, http.StatusBadRequest},
		{"failure:Invalid published date", "book", Book{Title: "the river", Author: Author{Id: 1}, Publication: "Penguin", PublishedDate: "1/7/1813"}, Book{}, http.StatusBadRequest},
		{"failure:duplicate id ", "book", Book{Title: "the river", Author: Author{Id: 1}, Publication: "Penguin", PublishedDate: "8/7/1900"}, Book{}, http.StatusBadRequest},
		{"failure:Invalid Id", "book", Book{Title: "little things", Author: Author{Id: 1}, Publication: "Penguin", PublishedDate: "3/7/1905"}, Book{}, http.StatusBadRequest},
	}
	for i, tc := range testCases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/"+tc.reqEndpoint, bytes.NewReader(body))
		postBook(w, req)
		res, _ := io.ReadAll(w.Result().Body)
		resBook := Book{}
		json.Unmarshal(res, &resBook)

		if reflect.DeepEqual(resBook, tc.expRes) {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestPostAuthor(t *testing.T) {
	testCases := []struct {
		desc        string
		reqEndpoint string
		reqBody     Author
		expRes      Author
		statusCode  int
	}{
		{"success case", "author", Author{FirstName: "Joey", LastName: "Paul", Dob: "01-08-2001", PenName: "Joe"}, Author{1, "Joey", "Paul", "01-08-2001", "Joe"}, http.StatusCreated},
		{"failure: first name is missing", "author", Author{FirstName: "", LastName: "Paul", Dob: "01-08-2001", PenName: "Joe"}, Author{}, http.StatusBadRequest},
		{"failure: last name is missing", "author", Author{FirstName: "", LastName: "", Dob: "01-08-2001", PenName: "Joe"}, Author{}, http.StatusBadRequest},
	}
	for i, tc := range testCases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/"+tc.reqEndpoint, bytes.NewReader(body))
		postAuthor(w, req)
		res, _ := io.ReadAll(w.Result().Body)
		resAuthor := Author{}
		json.Unmarshal(res, &resAuthor)

		if reflect.DeepEqual(resAuthor, tc.expRes) {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestPutBook(t *testing.T) {
	testcases := []struct {
		desc      string
		reqId     string
		reqBody   Book
		expRes    Book
		expStatus int
	}{
		{"successfully posted", "book", Book{Title: "the wall", Author: Author{Id: 1}, Publication: "Arihant", PublishedDate: "6/7/2017"}, Book{Title: "the wall", Author: Author{Id: 1}, Publication: "Arihant", PublishedDate: "6/7/2017"}, http.StatusCreated},
		{"successfully posted", "book", Book{Title: "city", Author: Author{Id: 1}, Publication: "Arihant", PublishedDate: "1/1/2001"}, Book{Title: "city", Author: Author{Id: 1}, Publication: "Arihant", PublishedDate: "1/1/2001"}, http.StatusCreated},
		{"failure : Invalid Publication", "book", Book{Title: "the wall", Author: Author{Id: 1}, Publication: "the sun", PublishedDate: "9/92012"}, Book{}, http.StatusBadRequest},
		{"failure:Invalid published date", "book", Book{Title: "the river", Author: Author{Id: 1}, Publication: "Penguin", PublishedDate: "1/7/1813"}, Book{}, http.StatusBadRequest},
		{"failure:duplicate id ", "book", Book{Title: "the river", Author: Author{Id: 1}, Publication: "Penguin", PublishedDate: "8/7/1900"}, Book{}, http.StatusBadRequest},
		{"failure:Invalid Id", "book", Book{Title: "little things", Author: Author{Id: 1}, Publication: "Penguin", PublishedDate: "3/7/1905"}, Book{}, http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/book/"+tc.reqId, bytes.NewReader(body))
		putBook(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
		res, _ := io.ReadAll(w.Result().Body)
		resBook := Book{}
		json.Unmarshal(res, &resBook)
		if resBook != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestPutAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		reqBody   Author
		expRes    Author
		expStatus int
	}{
		{"Valid details", Author{FirstName: "RD", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{1, "RD", "Sharma", "2/11/1989", "Sharma"}, http.StatusOK},
		{"InValid details", Author{FirstName: "", LastName: "Sharma", Dob: "2/11/1989", PenName: "Sharma"}, Author{}, http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		body, _ := json.Marshal(tc.reqBody)
		req := httptest.NewRequest(http.MethodPost, "localhost:8000/author", bytes.NewReader(body))
		putAuthor(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
		res, _ := io.ReadAll(w.Result().Body)
		resAuthor := Author{}
		json.Unmarshal(res, &resAuthor)
		if resAuthor != tc.expRes {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestDeleteBook(t *testing.T) {
	testcases := []struct {
		desc      string
		reqId     string
		expStatus int
	}{
		{"Valid Details", "1", http.StatusOK},
		{"Book does not exists", "100", http.StatusNotFound},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "localhost:8000/book/"+tc.reqId, nil)
		deleteBook(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}

func TestDeleteAuthor(t *testing.T) {
	testcases := []struct {
		desc      string
		reqId     string
		expStatus int
	}{
		{"Valid Details", "1", http.StatusOK},
		{"Author does not exists", "100", http.StatusBadRequest},
	}
	for i, tc := range testcases {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "localhost:8000/author/"+tc.reqId, nil)
		deleteAuthor(w, req)
		defer w.Result().Body.Close()

		if w.Result().StatusCode != tc.expStatus {
			t.Errorf("%v test failed %v", i, tc.desc)
		}
	}
}
