package folderfile

type FolderFile struct {
	File   string `json:"file"`
	Hash   string `json:"hash"`
	Result bool   `json:"result"`
	Error  string `json:"error"`
}

type FolderMap struct {
	Path   string       `json:"path"`
	Files  []FolderFile `json:"files"`
	Result bool         `json:"result"`
	Error  string       `json:"error"`
}
