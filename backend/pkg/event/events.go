package event

// BaseEvent 基础事件结构
type BaseEvent struct {
	Name string
}

// EventName 实现 Event 接口
func (e *BaseEvent) EventName() string {
	return e.Name
}

// 常见事件定义

// UserRegistered 用户注册事件
type UserRegistered struct {
	BaseEvent
	UserID   uint
	Username string
	Email    string
}

// NewUserRegistered 创建用户注册事件
func NewUserRegistered(userID uint, username, email string) *UserRegistered {
	return &UserRegistered{
		BaseEvent: BaseEvent{Name: "user.registered"},
		UserID:    userID,
		Username:  username,
		Email:     email,
	}
}

// UserLoggedIn 用户登录事件
type UserLoggedIn struct {
	BaseEvent
	UserID uint
	IP     string
}

// NewUserLoggedIn 创建用户登录事件
func NewUserLoggedIn(userID uint, ip string) *UserLoggedIn {
	return &UserLoggedIn{
		BaseEvent: BaseEvent{Name: "user.logged_in"},
		UserID:    userID,
		IP:        ip,
	}
}

// PostCreated 文章创建事件
type PostCreated struct {
	BaseEvent
	PostID   uint
	Title    string
	AuthorID uint
}

// NewPostCreated 创建文章创建事件
func NewPostCreated(postID uint, title string, authorID uint) *PostCreated {
	return &PostCreated{
		BaseEvent: BaseEvent{Name: "post.created"},
		PostID:    postID,
		Title:     title,
		AuthorID:  authorID,
	}
}

// PostPublished 文章发布事件
type PostPublished struct {
	BaseEvent
	PostID uint
	Title  string
	Slug   string
}

// NewPostPublished 创建文章发布事件
func NewPostPublished(postID uint, title, slug string) *PostPublished {
	return &PostPublished{
		BaseEvent: BaseEvent{Name: "post.published"},
		PostID:    postID,
		Title:     title,
		Slug:      slug,
	}
}

// CommentCreated 评论创建事件
type CommentCreated struct {
	BaseEvent
	CommentID uint
	PostID    uint
	Content   string
	AuthorID  uint
}

// NewCommentCreated 创建评论创建事件
func NewCommentCreated(commentID, postID uint, content string, authorID uint) *CommentCreated {
	return &CommentCreated{
		BaseEvent: BaseEvent{Name: "comment.created"},
		CommentID: commentID,
		PostID:    postID,
		Content:   content,
		AuthorID:  authorID,
	}
}

// FileUploaded 文件上传事件
type FileUploaded struct {
	BaseEvent
	FileID   uint
	FileName string
	FileSize int64
	UserID   uint
}

// NewFileUploaded 创建文件上传事件
func NewFileUploaded(fileID uint, fileName string, fileSize int64, userID uint) *FileUploaded {
	return &FileUploaded{
		BaseEvent: BaseEvent{Name: "file.uploaded"},
		FileID:    fileID,
		FileName:  fileName,
		FileSize:  fileSize,
		UserID:    userID,
	}
}

// OrderCreated 订单创建事件
type OrderCreated struct {
	BaseEvent
	OrderID    string
	UserID     uint
	TotalPrice float64
}

// NewOrderCreated 创建订单创建事件
func NewOrderCreated(orderID string, userID uint, totalPrice float64) *OrderCreated {
	return &OrderCreated{
		BaseEvent:  BaseEvent{Name: "order.created"},
		OrderID:    orderID,
		UserID:     userID,
		TotalPrice: totalPrice,
	}
}

// OrderPaid 订单支付事件
type OrderPaid struct {
	BaseEvent
	OrderID       string
	PaymentMethod string
	Amount        float64
}

// NewOrderPaid 创建订单支付事件
func NewOrderPaid(orderID, paymentMethod string, amount float64) *OrderPaid {
	return &OrderPaid{
		BaseEvent:     BaseEvent{Name: "order.paid"},
		OrderID:       orderID,
		PaymentMethod: paymentMethod,
		Amount:        amount,
	}
}
