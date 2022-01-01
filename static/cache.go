package static

import (
	"waho/comm"
	"waho/models"
)

// 缓存Key
const (
	// 用户token 存用户ID
	UserTokenKey = "user_token"
	// 用户列表 UserID是键
	UserKey = "user"
	// 当天手机验证码发送记录
	SendPhoneCodeCurrentDayLogKey = "send_phone_code_current_day_log"
	// 分类的缓存
	ClassifyMusterKey = "classify_muster"
	// 首页头banner的缓存
	BannerKey = "banner"
)

// 用户缓存
type UserCache struct {
	models.User
	CommKey string `json:"comm_key"` // 通讯加密key
	Token string `json:"token"`
}

// 发送手机验证码用户当前的缓存
type SendPhoneCodeCurrentDayLogCache struct {
	Phone string `json:"phone" validate:"phone"`
	Status int `json:"status"` // 是否成功 0:未发送、1:成功、2:失败
	Code string `json:"code"` // 验证码
	UserId int `json:"user_id"` // 发送的用户ID
	SendTime int64 `json:"send_time"` // 发送时间
}

// 分类
type Classify struct {
	Id int `json:"id" orm:"pk"`  //
	Title string `json:"title"`  //标题
	Subtitle string `json:"subtitle"`  //小标题
	Pid int `json:"pid"`  //上级分类
	Level int `json:"level"`  //分类级别一级分类、二级分类
	Image string `json:"image"`  //图片logo
	Type int `json:"type"`  //类型
	Link string `json:"link"`  //连接地址
}

// 分类树结构
type ClassifyMuster struct {
	Id int `json:"id" orm:"pk"`  //
	Title string `json:"title"`  //标题
	Subtitle string `json:"subtitle"`  //小标题
	Pid int `json:"pid"`  //上级分类
	Level int `json:"level"`  //分类级别一级分类、二级分类
	Image string `json:"image"`  //图片logo
	Type int `json:"type"`  //类型
	Link string `json:"link"`  //连接地址
	Children []Classify `json:"children"`
}

// 获取用户输出商品 包含是否有购买等
type Goods struct {
	Id                  int    `json:"id" orm:"pk"`           //
	Name                string `json:"name"`                  //标题
	OriginPrice         string `json:"origin_price"`          //原价格
	Price               string `json:"price"`                 //现价
	Spec                string `json:"spec"`                  //描述
	SpecTag             string `json:"spec_tag"`              //描述标签
	SmallImage          string `json:"small_image"`           //小图标
	TotalSales          string `json:"total_sales"`           //总销售额
	BuyLimit            int    `json:"buy_limit"`             //限买件数
	Stock               int    `json:"stock"`                 //库存
	MarkDiscount        int    `json:"mark_discount"`         //标记商品：折扣
	MarkNew             int    `json:"mark_new"`              //标记商品：新的
	MarkTag             string `json:"mark_tag"`              //标记商品：标签
	Status              int    `json:"status"`                //状态，0:正常、1:缺货、2:已下架、3:预售产品
	ClassifyTop         int    `json:"classify_top"`          //顶级分类，只有一个
	Classify            string `json:"classify"`              //所属类别[可多个,逗号分隔]，二级类别以下
	Type                int    `json:"type"`                  //类型
	Activity            string `json:"activity"`              //活动：逗号分隔
	PresaleDeliveryTime int    `json:"presale_delivery_time"` //预售发货时间
	IsBulk              int    `json:"is_bulk"`               //是否散装商品
	NetWeight           int    `json:"net_weight"`            //净重
	NetWeightUnit       string `json:"net_weight_unit"`       //净重单位
	IsInvoice           int    `json:"is_invoice"`            //是否可开发票
	AttributeTags       string `json:"attribute_tags"`        //属性标签，例[鲜活,现杀]
	DisableCouponsType  string `json:"disable_coupons_type"`  //禁用优惠券类型，例[all]
	Content             string `json:"content"`               //
}


func GetBannerKey(typeVal string, pk_id int) string {
	return typeVal + "_" + comm.ToSting(pk_id)
}