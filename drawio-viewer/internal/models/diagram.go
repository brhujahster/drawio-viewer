package models

type Diagram struct {
    ID      string `json:"id"`
    Name    string `json:"name"`
    XMLPath string `json:"xmlPath"`
    IsTemp  bool   `json:"isTemp"`
}