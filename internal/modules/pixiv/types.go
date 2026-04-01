package pixiv

type AjaxResponse[T any] struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Body    T      `json:"body"`
}

type illusDetailBody struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	CreateDate    string `json:"createDate"`
	UploadDate    string `json:"uploadDate"`
	Type          int    `json:"illustType"`
	PageCount     int    `json:"pageCount"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	BookmarkCount int    `json:"bookmarkCount"`
	LikeCount     int    `json:"likeCount"`
	ViewCount     int    `json:"viewCount"`
	UserID        string `json:"userId"`
	UserName      string `json:"userName"`
	UserAccount   string `json:"userAccount"`
	XRestrict     int    `json:"xRestrict"`
	Urls          struct {
		Mini     string `json:"mini"`
		Thumb    string `json:"thumb"`
		Small    string `json:"small"`
		Regular  string `json:"regular"`
		Original string `json:"original"`
	} `json:"urls"`
	Tags struct {
		Tags []struct {
			Tag         string `json:"tag"`
			Translation *struct {
				En string `json:"en"`
			} `json:"translation,omitempty"`
		} `json:"tags"`
	} `json:"tags"`
}

type IllustPagesBody []struct {
	URLs struct {
		ThumbMini string `json:"thumb_mini"`
		Small     string `json:"small"`
		Regular   string `json:"regular"`
		Original  string `json:"original"`
	} `json:"urls"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type UserFullBody struct {
	UserID     string `json:"userId"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ImageBig   string `json:"imageBig"`
	Premium    bool   `json:"premium"`
	IsFollowed bool   `json:"isFollowed"`
	IsBlocking bool   `json:"isBlocking"`
	Background struct {
		URL string `json:"url"`
	} `json:"background"`
	Social struct {
		Twitter interface{} `json:"twitter"`
	} `json:"social"`
}

type UserProfileAllBody struct {
	Illusts map[string]struct{} `json:"illusts"`
	Manga   map[string]struct{} `json:"manga"`
}

type SearchArtworksBody struct {
	Illust struct {
		Data []struct {
			ID            string   `json:"id"`
			Title         string   `json:"title"`
			Alt           string   `json:"alt"`
			URL           string   `json:"url"`
			UserID        string   `json:"userId"`
			UserName      string   `json:"userName"`
			PageCount     int      `json:"pageCount"`
			BookmarkCount int      `json:"bookmarkCount"`
			Width         int      `json:"width"`
			Height        int      `json:"height"`
			IllustType    int      `json:"illustType"`
			XRestrict     int      `json:"xRestrict"`
			CreateDate    string   `json:"createDate"`
			Description   string   `json:"description"`
			Tags          []string `json:"tags"`
		} `json:"data"`
		Total string `json:"total"`
	} `json:"illust"`
}
