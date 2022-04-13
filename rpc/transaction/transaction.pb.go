// Code generated by protoc-gen-go. DO NOT EDIT.
// source: transaction.proto

package transaction

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ProxyOrderReq_DFB struct {
	Req                  *ProxyPayOrderRequest `protobuf:"bytes,1,opt,name=req,proto3" json:"req,omitempty"`
	Rate                 *CorrespondMerChnRate `protobuf:"bytes,2,opt,name=rate,proto3" json:"rate,omitempty"`
	BalanceType          string                `protobuf:"bytes,3,opt,name=balanceType,proto3" json:"balanceType,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *ProxyOrderReq_DFB) Reset()         { *m = ProxyOrderReq_DFB{} }
func (m *ProxyOrderReq_DFB) String() string { return proto.CompactTextString(m) }
func (*ProxyOrderReq_DFB) ProtoMessage()    {}
func (*ProxyOrderReq_DFB) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{0}
}

func (m *ProxyOrderReq_DFB) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyOrderReq_DFB.Unmarshal(m, b)
}
func (m *ProxyOrderReq_DFB) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyOrderReq_DFB.Marshal(b, m, deterministic)
}
func (m *ProxyOrderReq_DFB) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyOrderReq_DFB.Merge(m, src)
}
func (m *ProxyOrderReq_DFB) XXX_Size() int {
	return xxx_messageInfo_ProxyOrderReq_DFB.Size(m)
}
func (m *ProxyOrderReq_DFB) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyOrderReq_DFB.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyOrderReq_DFB proto.InternalMessageInfo

func (m *ProxyOrderReq_DFB) GetReq() *ProxyPayOrderRequest {
	if m != nil {
		return m.Req
	}
	return nil
}

func (m *ProxyOrderReq_DFB) GetRate() *CorrespondMerChnRate {
	if m != nil {
		return m.Rate
	}
	return nil
}

func (m *ProxyOrderReq_DFB) GetBalanceType() string {
	if m != nil {
		return m.BalanceType
	}
	return ""
}

type ProxyOrderResp_DFB struct {
	ProxyOrderNo         string   `protobuf:"bytes,1,opt,name=ProxyOrderNo,proto3" json:"ProxyOrderNo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProxyOrderResp_DFB) Reset()         { *m = ProxyOrderResp_DFB{} }
func (m *ProxyOrderResp_DFB) String() string { return proto.CompactTextString(m) }
func (*ProxyOrderResp_DFB) ProtoMessage()    {}
func (*ProxyOrderResp_DFB) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{1}
}

func (m *ProxyOrderResp_DFB) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyOrderResp_DFB.Unmarshal(m, b)
}
func (m *ProxyOrderResp_DFB) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyOrderResp_DFB.Marshal(b, m, deterministic)
}
func (m *ProxyOrderResp_DFB) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyOrderResp_DFB.Merge(m, src)
}
func (m *ProxyOrderResp_DFB) XXX_Size() int {
	return xxx_messageInfo_ProxyOrderResp_DFB.Size(m)
}
func (m *ProxyOrderResp_DFB) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyOrderResp_DFB.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyOrderResp_DFB proto.InternalMessageInfo

func (m *ProxyOrderResp_DFB) GetProxyOrderNo() string {
	if m != nil {
		return m.ProxyOrderNo
	}
	return ""
}

type ProxyPayOrderRequest struct {
	AccessType           string   `protobuf:"bytes,1,opt,name=AccessType,proto3" json:"AccessType,omitempty"`
	MerchantId           string   `protobuf:"bytes,2,opt,name=MerchantId,proto3" json:"MerchantId,omitempty"`
	Sign                 string   `protobuf:"bytes,3,opt,name=Sign,proto3" json:"Sign,omitempty"`
	NotifyUrl            string   `protobuf:"bytes,4,opt,name=NotifyUrl,proto3" json:"NotifyUrl,omitempty"`
	Language             string   `protobuf:"bytes,5,opt,name=Language,proto3" json:"Language,omitempty"`
	OrderNo              string   `protobuf:"bytes,6,opt,name=OrderNo,proto3" json:"OrderNo,omitempty"`
	BankId               string   `protobuf:"bytes,7,opt,name=BankId,proto3" json:"BankId,omitempty"`
	BankName             string   `protobuf:"bytes,8,opt,name=BankName,proto3" json:"BankName,omitempty"`
	BankProvince         string   `protobuf:"bytes,9,opt,name=BankProvince,proto3" json:"BankProvince,omitempty"`
	BankCity             string   `protobuf:"bytes,10,opt,name=BankCity,proto3" json:"BankCity,omitempty"`
	BranchName           string   `protobuf:"bytes,11,opt,name=BranchName,proto3" json:"BranchName,omitempty"`
	BankNo               string   `protobuf:"bytes,12,opt,name=BankNo,proto3" json:"BankNo,omitempty"`
	OrderAmount          float64  `protobuf:"fixed64,13,opt,name=OrderAmount,proto3" json:"OrderAmount,omitempty"`
	DefrayName           string   `protobuf:"bytes,14,opt,name=DefrayName,proto3" json:"DefrayName,omitempty"`
	DefrayId             string   `protobuf:"bytes,15,opt,name=DefrayId,proto3" json:"DefrayId,omitempty"`
	DefrayMobile         string   `protobuf:"bytes,16,opt,name=DefrayMobile,proto3" json:"DefrayMobile,omitempty"`
	DefrayEmail          string   `protobuf:"bytes,17,opt,name=DefrayEmail,proto3" json:"DefrayEmail,omitempty"`
	Currency             string   `protobuf:"bytes,18,opt,name=Currency,proto3" json:"Currency,omitempty"`
	PayTypeSubNo         string   `protobuf:"bytes,19,opt,name=PayTypeSubNo,proto3" json:"PayTypeSubNo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProxyPayOrderRequest) Reset()         { *m = ProxyPayOrderRequest{} }
func (m *ProxyPayOrderRequest) String() string { return proto.CompactTextString(m) }
func (*ProxyPayOrderRequest) ProtoMessage()    {}
func (*ProxyPayOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{2}
}

func (m *ProxyPayOrderRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProxyPayOrderRequest.Unmarshal(m, b)
}
func (m *ProxyPayOrderRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProxyPayOrderRequest.Marshal(b, m, deterministic)
}
func (m *ProxyPayOrderRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProxyPayOrderRequest.Merge(m, src)
}
func (m *ProxyPayOrderRequest) XXX_Size() int {
	return xxx_messageInfo_ProxyPayOrderRequest.Size(m)
}
func (m *ProxyPayOrderRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProxyPayOrderRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProxyPayOrderRequest proto.InternalMessageInfo

func (m *ProxyPayOrderRequest) GetAccessType() string {
	if m != nil {
		return m.AccessType
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetMerchantId() string {
	if m != nil {
		return m.MerchantId
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetSign() string {
	if m != nil {
		return m.Sign
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetNotifyUrl() string {
	if m != nil {
		return m.NotifyUrl
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetLanguage() string {
	if m != nil {
		return m.Language
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetOrderNo() string {
	if m != nil {
		return m.OrderNo
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetBankId() string {
	if m != nil {
		return m.BankId
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetBankName() string {
	if m != nil {
		return m.BankName
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetBankProvince() string {
	if m != nil {
		return m.BankProvince
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetBankCity() string {
	if m != nil {
		return m.BankCity
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetBranchName() string {
	if m != nil {
		return m.BranchName
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetBankNo() string {
	if m != nil {
		return m.BankNo
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetOrderAmount() float64 {
	if m != nil {
		return m.OrderAmount
	}
	return 0
}

func (m *ProxyPayOrderRequest) GetDefrayName() string {
	if m != nil {
		return m.DefrayName
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetDefrayId() string {
	if m != nil {
		return m.DefrayId
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetDefrayMobile() string {
	if m != nil {
		return m.DefrayMobile
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetDefrayEmail() string {
	if m != nil {
		return m.DefrayEmail
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetCurrency() string {
	if m != nil {
		return m.Currency
	}
	return ""
}

func (m *ProxyPayOrderRequest) GetPayTypeSubNo() string {
	if m != nil {
		return m.PayTypeSubNo
	}
	return ""
}

type CorrespondMerChnRate struct {
	MerchantCode         string   `protobuf:"bytes,1,opt,name=MerchantCode,proto3" json:"MerchantCode,omitempty"`
	ChannelPayTypesCode  string   `protobuf:"bytes,2,opt,name=ChannelPayTypesCode,proto3" json:"ChannelPayTypesCode,omitempty"`
	ChannelCode          string   `protobuf:"bytes,3,opt,name=ChannelCode,proto3" json:"ChannelCode,omitempty"`
	PayTypeCode          string   `protobuf:"bytes,4,opt,name=PayTypeCode,proto3" json:"PayTypeCode,omitempty"`
	PayTypeCodeNum       string   `protobuf:"bytes,5,opt,name=PayTypeCodeNum,proto3" json:"PayTypeCodeNum,omitempty"`
	Designation          string   `protobuf:"bytes,6,opt,name=Designation,proto3" json:"Designation,omitempty"`
	DesignationNo        string   `protobuf:"bytes,7,opt,name=DesignationNo,proto3" json:"DesignationNo,omitempty"`
	Fee                  float64  `protobuf:"fixed64,8,opt,name=Fee,proto3" json:"Fee,omitempty"`
	HandlingFee          float64  `protobuf:"fixed64,9,opt,name=HandlingFee,proto3" json:"HandlingFee,omitempty"`
	ChFee                float64  `protobuf:"fixed64,10,opt,name=ChFee,proto3" json:"ChFee,omitempty"`
	ChHandlingFee        float64  `protobuf:"fixed64,11,opt,name=ChHandlingFee,proto3" json:"ChHandlingFee,omitempty"`
	SingleMinCharge      float64  `protobuf:"fixed64,12,opt,name=SingleMinCharge,proto3" json:"SingleMinCharge,omitempty"`
	SingleMaxCharge      float64  `protobuf:"fixed64,13,opt,name=SingleMaxCharge,proto3" json:"SingleMaxCharge,omitempty"`
	CurrencyCode         string   `protobuf:"bytes,14,opt,name=CurrencyCode,proto3" json:"CurrencyCode,omitempty"`
	ApiUrl               string   `protobuf:"bytes,15,opt,name=ApiUrl,proto3" json:"ApiUrl,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CorrespondMerChnRate) Reset()         { *m = CorrespondMerChnRate{} }
func (m *CorrespondMerChnRate) String() string { return proto.CompactTextString(m) }
func (*CorrespondMerChnRate) ProtoMessage()    {}
func (*CorrespondMerChnRate) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{3}
}

func (m *CorrespondMerChnRate) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CorrespondMerChnRate.Unmarshal(m, b)
}
func (m *CorrespondMerChnRate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CorrespondMerChnRate.Marshal(b, m, deterministic)
}
func (m *CorrespondMerChnRate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CorrespondMerChnRate.Merge(m, src)
}
func (m *CorrespondMerChnRate) XXX_Size() int {
	return xxx_messageInfo_CorrespondMerChnRate.Size(m)
}
func (m *CorrespondMerChnRate) XXX_DiscardUnknown() {
	xxx_messageInfo_CorrespondMerChnRate.DiscardUnknown(m)
}

var xxx_messageInfo_CorrespondMerChnRate proto.InternalMessageInfo

func (m *CorrespondMerChnRate) GetMerchantCode() string {
	if m != nil {
		return m.MerchantCode
	}
	return ""
}

func (m *CorrespondMerChnRate) GetChannelPayTypesCode() string {
	if m != nil {
		return m.ChannelPayTypesCode
	}
	return ""
}

func (m *CorrespondMerChnRate) GetChannelCode() string {
	if m != nil {
		return m.ChannelCode
	}
	return ""
}

func (m *CorrespondMerChnRate) GetPayTypeCode() string {
	if m != nil {
		return m.PayTypeCode
	}
	return ""
}

func (m *CorrespondMerChnRate) GetPayTypeCodeNum() string {
	if m != nil {
		return m.PayTypeCodeNum
	}
	return ""
}

func (m *CorrespondMerChnRate) GetDesignation() string {
	if m != nil {
		return m.Designation
	}
	return ""
}

func (m *CorrespondMerChnRate) GetDesignationNo() string {
	if m != nil {
		return m.DesignationNo
	}
	return ""
}

func (m *CorrespondMerChnRate) GetFee() float64 {
	if m != nil {
		return m.Fee
	}
	return 0
}

func (m *CorrespondMerChnRate) GetHandlingFee() float64 {
	if m != nil {
		return m.HandlingFee
	}
	return 0
}

func (m *CorrespondMerChnRate) GetChFee() float64 {
	if m != nil {
		return m.ChFee
	}
	return 0
}

func (m *CorrespondMerChnRate) GetChHandlingFee() float64 {
	if m != nil {
		return m.ChHandlingFee
	}
	return 0
}

func (m *CorrespondMerChnRate) GetSingleMinCharge() float64 {
	if m != nil {
		return m.SingleMinCharge
	}
	return 0
}

func (m *CorrespondMerChnRate) GetSingleMaxCharge() float64 {
	if m != nil {
		return m.SingleMaxCharge
	}
	return 0
}

func (m *CorrespondMerChnRate) GetCurrencyCode() string {
	if m != nil {
		return m.CurrencyCode
	}
	return ""
}

func (m *CorrespondMerChnRate) GetApiUrl() string {
	if m != nil {
		return m.ApiUrl
	}
	return ""
}

type PayOrderRequest struct {
	PayOrder             *PayOrder             `protobuf:"bytes,1,opt,name=payOrder,proto3" json:"payOrder,omitempty"`
	Rate                 *CorrespondMerChnRate `protobuf:"bytes,2,opt,name=Rate,proto3" json:"Rate,omitempty"`
	OrderNo              string                `protobuf:"bytes,3,opt,name=OrderNo,proto3" json:"OrderNo,omitempty"`
	ChannelOrderNo       string                `protobuf:"bytes,4,opt,name=ChannelOrderNo,proto3" json:"ChannelOrderNo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}              `json:"-"`
	XXX_unrecognized     []byte                `json:"-"`
	XXX_sizecache        int32                 `json:"-"`
}

func (m *PayOrderRequest) Reset()         { *m = PayOrderRequest{} }
func (m *PayOrderRequest) String() string { return proto.CompactTextString(m) }
func (*PayOrderRequest) ProtoMessage()    {}
func (*PayOrderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{4}
}

func (m *PayOrderRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PayOrderRequest.Unmarshal(m, b)
}
func (m *PayOrderRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PayOrderRequest.Marshal(b, m, deterministic)
}
func (m *PayOrderRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayOrderRequest.Merge(m, src)
}
func (m *PayOrderRequest) XXX_Size() int {
	return xxx_messageInfo_PayOrderRequest.Size(m)
}
func (m *PayOrderRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PayOrderRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PayOrderRequest proto.InternalMessageInfo

func (m *PayOrderRequest) GetPayOrder() *PayOrder {
	if m != nil {
		return m.PayOrder
	}
	return nil
}

func (m *PayOrderRequest) GetRate() *CorrespondMerChnRate {
	if m != nil {
		return m.Rate
	}
	return nil
}

func (m *PayOrderRequest) GetOrderNo() string {
	if m != nil {
		return m.OrderNo
	}
	return ""
}

func (m *PayOrderRequest) GetChannelOrderNo() string {
	if m != nil {
		return m.ChannelOrderNo
	}
	return ""
}

type PayOrderResponse struct {
	PayOrderNo           string   `protobuf:"bytes,1,opt,name=PayOrderNo,proto3" json:"PayOrderNo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PayOrderResponse) Reset()         { *m = PayOrderResponse{} }
func (m *PayOrderResponse) String() string { return proto.CompactTextString(m) }
func (*PayOrderResponse) ProtoMessage()    {}
func (*PayOrderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{5}
}

func (m *PayOrderResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PayOrderResponse.Unmarshal(m, b)
}
func (m *PayOrderResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PayOrderResponse.Marshal(b, m, deterministic)
}
func (m *PayOrderResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayOrderResponse.Merge(m, src)
}
func (m *PayOrderResponse) XXX_Size() int {
	return xxx_messageInfo_PayOrderResponse.Size(m)
}
func (m *PayOrderResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PayOrderResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PayOrderResponse proto.InternalMessageInfo

func (m *PayOrderResponse) GetPayOrderNo() string {
	if m != nil {
		return m.PayOrderNo
	}
	return ""
}

type PayOrder struct {
	AccessType           string   `protobuf:"bytes,1,opt,name=AccessType,proto3" json:"AccessType,omitempty"`
	MerchantId           string   `protobuf:"bytes,2,opt,name=MerchantId,proto3" json:"MerchantId,omitempty"`
	NotifyUrl            string   `protobuf:"bytes,3,opt,name=NotifyUrl,proto3" json:"NotifyUrl,omitempty"`
	PageUrl              string   `protobuf:"bytes,4,opt,name=PageUrl,proto3" json:"PageUrl,omitempty"`
	Language             string   `protobuf:"bytes,5,opt,name=Language,proto3" json:"Language,omitempty"`
	Sign                 string   `protobuf:"bytes,6,opt,name=Sign,proto3" json:"Sign,omitempty"`
	JumpType             string   `protobuf:"bytes,7,opt,name=JumpType,proto3" json:"JumpType,omitempty"`
	LoginType            string   `protobuf:"bytes,8,opt,name=LoginType,proto3" json:"LoginType,omitempty"`
	OrderNo              string   `protobuf:"bytes,9,opt,name=OrderNo,proto3" json:"OrderNo,omitempty"`
	OrderAmount          string   `protobuf:"bytes,10,opt,name=OrderAmount,proto3" json:"OrderAmount,omitempty"`
	Currency             string   `protobuf:"bytes,11,opt,name=Currency,proto3" json:"Currency,omitempty"`
	PayType              string   `protobuf:"bytes,12,opt,name=PayType,proto3" json:"PayType,omitempty"`
	OrderTime            string   `protobuf:"bytes,13,opt,name=OrderTime,proto3" json:"OrderTime,omitempty"`
	OrderName            string   `protobuf:"bytes,14,opt,name=OrderName,proto3" json:"OrderName,omitempty"`
	BankCode             string   `protobuf:"bytes,15,opt,name=BankCode,proto3" json:"BankCode,omitempty"`
	Phone                string   `protobuf:"bytes,16,opt,name=Phone,proto3" json:"Phone,omitempty"`
	PayTypeNo            string   `protobuf:"bytes,17,opt,name=PayTypeNo,proto3" json:"PayTypeNo,omitempty"`
	UserId               string   `protobuf:"bytes,18,opt,name=UserId,proto3" json:"UserId,omitempty"`
	MerchantLevel        string   `protobuf:"bytes,19,opt,name=MerchantLevel,proto3" json:"MerchantLevel,omitempty"`
	UserIp               string   `protobuf:"bytes,20,opt,name=UserIp,proto3" json:"UserIp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PayOrder) Reset()         { *m = PayOrder{} }
func (m *PayOrder) String() string { return proto.CompactTextString(m) }
func (*PayOrder) ProtoMessage()    {}
func (*PayOrder) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cc4e03d2c28c490, []int{6}
}

func (m *PayOrder) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PayOrder.Unmarshal(m, b)
}
func (m *PayOrder) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PayOrder.Marshal(b, m, deterministic)
}
func (m *PayOrder) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PayOrder.Merge(m, src)
}
func (m *PayOrder) XXX_Size() int {
	return xxx_messageInfo_PayOrder.Size(m)
}
func (m *PayOrder) XXX_DiscardUnknown() {
	xxx_messageInfo_PayOrder.DiscardUnknown(m)
}

var xxx_messageInfo_PayOrder proto.InternalMessageInfo

func (m *PayOrder) GetAccessType() string {
	if m != nil {
		return m.AccessType
	}
	return ""
}

func (m *PayOrder) GetMerchantId() string {
	if m != nil {
		return m.MerchantId
	}
	return ""
}

func (m *PayOrder) GetNotifyUrl() string {
	if m != nil {
		return m.NotifyUrl
	}
	return ""
}

func (m *PayOrder) GetPageUrl() string {
	if m != nil {
		return m.PageUrl
	}
	return ""
}

func (m *PayOrder) GetLanguage() string {
	if m != nil {
		return m.Language
	}
	return ""
}

func (m *PayOrder) GetSign() string {
	if m != nil {
		return m.Sign
	}
	return ""
}

func (m *PayOrder) GetJumpType() string {
	if m != nil {
		return m.JumpType
	}
	return ""
}

func (m *PayOrder) GetLoginType() string {
	if m != nil {
		return m.LoginType
	}
	return ""
}

func (m *PayOrder) GetOrderNo() string {
	if m != nil {
		return m.OrderNo
	}
	return ""
}

func (m *PayOrder) GetOrderAmount() string {
	if m != nil {
		return m.OrderAmount
	}
	return ""
}

func (m *PayOrder) GetCurrency() string {
	if m != nil {
		return m.Currency
	}
	return ""
}

func (m *PayOrder) GetPayType() string {
	if m != nil {
		return m.PayType
	}
	return ""
}

func (m *PayOrder) GetOrderTime() string {
	if m != nil {
		return m.OrderTime
	}
	return ""
}

func (m *PayOrder) GetOrderName() string {
	if m != nil {
		return m.OrderName
	}
	return ""
}

func (m *PayOrder) GetBankCode() string {
	if m != nil {
		return m.BankCode
	}
	return ""
}

func (m *PayOrder) GetPhone() string {
	if m != nil {
		return m.Phone
	}
	return ""
}

func (m *PayOrder) GetPayTypeNo() string {
	if m != nil {
		return m.PayTypeNo
	}
	return ""
}

func (m *PayOrder) GetUserId() string {
	if m != nil {
		return m.UserId
	}
	return ""
}

func (m *PayOrder) GetMerchantLevel() string {
	if m != nil {
		return m.MerchantLevel
	}
	return ""
}

func (m *PayOrder) GetUserIp() string {
	if m != nil {
		return m.UserIp
	}
	return ""
}

func init() {
	proto.RegisterType((*ProxyOrderReq_DFB)(nil), "transaction.ProxyOrderReq_DFB")
	proto.RegisterType((*ProxyOrderResp_DFB)(nil), "transaction.ProxyOrderResp_DFB")
	proto.RegisterType((*ProxyPayOrderRequest)(nil), "transaction.ProxyPayOrderRequest")
	proto.RegisterType((*CorrespondMerChnRate)(nil), "transaction.CorrespondMerChnRate")
	proto.RegisterType((*PayOrderRequest)(nil), "transaction.PayOrderRequest")
	proto.RegisterType((*PayOrderResponse)(nil), "transaction.PayOrderResponse")
	proto.RegisterType((*PayOrder)(nil), "transaction.PayOrder")
}

func init() { proto.RegisterFile("transaction.proto", fileDescriptor_2cc4e03d2c28c490) }

var fileDescriptor_2cc4e03d2c28c490 = []byte{
	// 910 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x56, 0xdd, 0x6e, 0xe3, 0x44,
	0x14, 0x96, 0x49, 0xda, 0x26, 0x27, 0xdb, 0xa6, 0x9d, 0x0d, 0x68, 0x54, 0x2d, 0x4b, 0x89, 0xd0,
	0xaa, 0x57, 0x05, 0xba, 0x42, 0xe2, 0xb6, 0xf5, 0xb2, 0xa2, 0xa8, 0x35, 0xc1, 0xed, 0x0a, 0x89,
	0x1b, 0x34, 0xb1, 0x67, 0x6d, 0x0b, 0x67, 0xc6, 0x3b, 0x76, 0x56, 0x9b, 0xb7, 0x81, 0xc7, 0xe0,
	0x02, 0x89, 0x27, 0xe0, 0x99, 0xd0, 0x99, 0x1f, 0x7b, 0x1c, 0xb2, 0x48, 0x88, 0xbb, 0xf9, 0xbe,
	0x73, 0xe6, 0xcc, 0x99, 0xc9, 0x77, 0xbe, 0x18, 0x4e, 0x1a, 0xc5, 0x44, 0xcd, 0x92, 0xa6, 0x90,
	0xe2, 0xa2, 0x52, 0xb2, 0x91, 0x64, 0xe2, 0x51, 0xf3, 0xdf, 0x02, 0x38, 0x59, 0x28, 0xf9, 0x6e,
	0xf3, 0xbd, 0x4a, 0xb9, 0x8a, 0xf9, 0x9b, 0x9f, 0x5f, 0xbc, 0xbc, 0x26, 0xcf, 0x61, 0xa0, 0xf8,
	0x1b, 0x1a, 0x9c, 0x05, 0xe7, 0x93, 0xcb, 0x4f, 0x2f, 0xfc, 0x1a, 0x3a, 0x79, 0xc1, 0xda, 0xfc,
	0x35, 0xaf, 0x9b, 0x18, 0xb3, 0xc9, 0x57, 0x30, 0x54, 0xac, 0xe1, 0xf4, 0x83, 0x1d, 0xbb, 0x42,
	0xa9, 0x14, 0xaf, 0x2b, 0x29, 0xd2, 0x3b, 0xae, 0xc2, 0x5c, 0xc4, 0xac, 0xe1, 0xb1, 0x4e, 0x27,
	0x67, 0x30, 0x59, 0xb2, 0x92, 0x89, 0x84, 0x3f, 0x6c, 0x2a, 0x4e, 0x07, 0x67, 0xc1, 0xf9, 0x38,
	0xf6, 0xa9, 0xf9, 0xd7, 0x40, 0xfc, 0x16, 0xeb, 0x4a, 0xf7, 0x38, 0x87, 0x47, 0x1d, 0x1b, 0x49,
	0xdd, 0xec, 0x38, 0xee, 0x71, 0xf3, 0xbf, 0x86, 0x30, 0xdb, 0xd5, 0x30, 0x79, 0x0a, 0x70, 0x95,
	0x24, 0xbc, 0xae, 0xf5, 0x99, 0x66, 0xab, 0xc7, 0x60, 0xfc, 0x8e, 0xab, 0x24, 0x67, 0xa2, 0xb9,
	0x49, 0xf5, 0x8d, 0xc6, 0xb1, 0xc7, 0x10, 0x02, 0xc3, 0xfb, 0x22, 0x13, 0xb6, 0x5b, 0xbd, 0x26,
	0x4f, 0x60, 0x1c, 0xc9, 0xa6, 0x78, 0xbd, 0x79, 0xa5, 0x4a, 0x3a, 0xd4, 0x81, 0x8e, 0x20, 0xa7,
	0x30, 0xba, 0x65, 0x22, 0x5b, 0xb3, 0x8c, 0xd3, 0x3d, 0x1d, 0x6c, 0x31, 0xa1, 0x70, 0xe0, 0x6e,
	0xb1, 0xaf, 0x43, 0x0e, 0x92, 0x8f, 0x60, 0xff, 0x9a, 0x89, 0x5f, 0x6e, 0x52, 0x7a, 0xa0, 0x03,
	0x16, 0x61, 0x35, 0x5c, 0x45, 0x6c, 0xc5, 0xe9, 0xc8, 0x54, 0x73, 0x18, 0x1f, 0x06, 0xd7, 0x0b,
	0x25, 0xdf, 0x16, 0x22, 0xe1, 0x74, 0x6c, 0x1e, 0xc6, 0xe7, 0xdc, 0xfe, 0xb0, 0x68, 0x36, 0x14,
	0xba, 0xfd, 0x88, 0xf1, 0xee, 0xd7, 0x8a, 0x89, 0x24, 0xd7, 0xd5, 0x27, 0xe6, 0xee, 0x1d, 0xe3,
	0x7a, 0x8a, 0x24, 0x7d, 0xd4, 0xf5, 0x14, 0x49, 0xfc, 0x21, 0x75, 0xdb, 0x57, 0x2b, 0xb9, 0x16,
	0x0d, 0x3d, 0x3c, 0x0b, 0xce, 0x83, 0xd8, 0xa7, 0xb0, 0xf2, 0x0b, 0xfe, 0x5a, 0xb1, 0x8d, 0xae,
	0x7c, 0x64, 0x2a, 0x77, 0x0c, 0x76, 0x65, 0xd0, 0x4d, 0x4a, 0xa7, 0xa6, 0x2b, 0x87, 0xf1, 0x56,
	0x66, 0x7d, 0x27, 0x97, 0x45, 0xc9, 0xe9, 0xb1, 0xb9, 0x95, 0xcf, 0x61, 0x07, 0x06, 0x7f, 0xb3,
	0x62, 0x45, 0x49, 0x4f, 0x8c, 0x94, 0x3c, 0x0a, 0x4f, 0x08, 0xd7, 0x4a, 0x71, 0x91, 0x6c, 0x28,
	0x31, 0x27, 0x38, 0xac, 0x05, 0xc5, 0x36, 0xf8, 0xf3, 0xdf, 0xaf, 0x97, 0x91, 0xa4, 0x8f, 0xad,
	0xa0, 0x3c, 0x6e, 0xfe, 0xeb, 0x10, 0x66, 0xbb, 0xb4, 0x8c, 0x9b, 0x9d, 0x3c, 0x42, 0x99, 0x3a,
	0x49, 0xf5, 0x38, 0xf2, 0x05, 0x3c, 0x0e, 0x73, 0x26, 0x04, 0x2f, 0x6d, 0xcd, 0x5a, 0xa7, 0x1a,
	0x75, 0xed, 0x0a, 0xe1, 0x85, 0x2c, 0xad, 0x33, 0xed, 0x6c, 0x78, 0x14, 0x66, 0xd8, 0x1d, 0x3a,
	0xc3, 0xc8, 0xce, 0xa7, 0xc8, 0x33, 0x38, 0xf2, 0x60, 0xb4, 0x5e, 0x59, 0xf9, 0x6d, 0xb1, 0xe6,
	0xf1, 0xea, 0x22, 0x13, 0x0c, 0x27, 0xd6, 0x0a, 0xd1, 0xa7, 0xc8, 0x67, 0x70, 0xe8, 0xc1, 0x48,
	0x5a, 0x4d, 0xf6, 0x49, 0x72, 0x0c, 0x83, 0x97, 0xdc, 0xa8, 0x32, 0x88, 0x71, 0x89, 0x95, 0xbf,
	0x65, 0x22, 0x2d, 0x0b, 0x91, 0x61, 0x64, 0x6c, 0x84, 0xe1, 0x51, 0x64, 0x06, 0x7b, 0x61, 0x8e,
	0x31, 0xd0, 0x31, 0x03, 0xf0, 0xbc, 0x30, 0xf7, 0x77, 0x4e, 0x74, 0xb4, 0x4f, 0x92, 0x73, 0x98,
	0xde, 0x17, 0x22, 0x2b, 0xf9, 0x5d, 0x21, 0xc2, 0x9c, 0xa9, 0x8c, 0x6b, 0x5d, 0x06, 0xf1, 0x36,
	0xed, 0x65, 0xb2, 0x77, 0x36, 0xf3, 0xb0, 0x97, 0xe9, 0x68, 0xfc, 0x35, 0x9d, 0x2c, 0xf4, 0xb3,
	0x1a, 0xa9, 0xf6, 0x38, 0x1c, 0x83, 0xab, 0xaa, 0xc0, 0x59, 0x37, 0x52, 0xb5, 0x68, 0xfe, 0x47,
	0x00, 0xd3, 0x6d, 0xbb, 0xf9, 0x12, 0x46, 0x95, 0xa5, 0xac, 0xa9, 0x7e, 0xd8, 0x37, 0x55, 0x97,
	0xdf, 0xa6, 0xa1, 0x9b, 0xc6, 0xff, 0xcd, 0x4d, 0xb5, 0x0e, 0x3d, 0x2b, 0x19, 0xf4, 0xad, 0xe4,
	0x19, 0x1c, 0x59, 0xe1, 0xb8, 0x04, 0x23, 0x96, 0x2d, 0x76, 0x7e, 0x09, 0xc7, 0x5d, 0xfb, 0x78,
	0x48, 0xad, 0xed, 0xd0, 0x71, 0xad, 0xd3, 0x7a, 0xcc, 0xfc, 0xcf, 0x21, 0x8c, 0x1c, 0xfc, 0xdf,
	0xde, 0xda, 0xf3, 0xd1, 0xc1, 0xb6, 0x8f, 0x52, 0x38, 0x58, 0xb0, 0x8c, 0x77, 0x1e, 0xeb, 0xe0,
	0xbf, 0x3a, 0xac, 0xf3, 0xeb, 0x7d, 0xcf, 0xaf, 0x4f, 0x61, 0xf4, 0xdd, 0x7a, 0x55, 0xe9, 0x2e,
	0x8d, 0x92, 0x5b, 0x8c, 0x3d, 0xdc, 0xca, 0xac, 0x10, 0x3a, 0x68, 0x0c, 0xb6, 0x23, 0xfc, 0x47,
	0x1e, 0xf7, 0x1f, 0x79, 0xcb, 0x03, 0x8d, 0xb5, 0xf6, 0x3c, 0xd0, 0x77, 0xa0, 0xc9, 0x96, 0x03,
	0xe9, 0xbb, 0xe9, 0xa1, 0xb4, 0xd6, 0xea, 0x20, 0xf6, 0xa3, 0x8b, 0x3c, 0x14, 0x2b, 0x23, 0xda,
	0x71, 0xdc, 0x11, 0x6d, 0xd4, 0xb3, 0xd5, 0x8e, 0x68, 0xbd, 0x1e, 0x85, 0x3c, 0xf5, 0xbc, 0x1e,
	0x45, 0x3c, 0x83, 0xbd, 0x45, 0x2e, 0x85, 0xb3, 0x53, 0x03, 0xb0, 0x9e, 0x3d, 0x38, 0x92, 0xd6,
	0x45, 0x3b, 0x02, 0x85, 0xff, 0xaa, 0xe6, 0xea, 0x26, 0xb5, 0x0e, 0x6a, 0x11, 0x8e, 0xab, 0xfb,
	0x15, 0x6f, 0xf9, 0x5b, 0x5e, 0x5a, 0x03, 0xed, 0x93, 0xed, 0xee, 0x8a, 0xce, 0xbc, 0xdd, 0xd5,
	0xe5, 0xef, 0x01, 0x4c, 0x1e, 0x3a, 0x8d, 0x93, 0x1f, 0xed, 0x3f, 0xb7, 0xb9, 0xa5, 0x62, 0xc2,
	0xf2, 0x4f, 0xff, 0xf9, 0x35, 0xe2, 0x7f, 0xba, 0x9c, 0x7e, 0xf2, 0xde, 0xb8, 0xfd, 0x6e, 0xf8,
	0x01, 0x88, 0x93, 0xaa, 0x57, 0xf6, 0xc9, 0xee, 0x79, 0x34, 0xf3, 0x7b, 0xfa, 0xf1, 0x7b, 0xa2,
	0x66, 0x3c, 0xae, 0xa7, 0x3f, 0x1d, 0x5e, 0x7c, 0xee, 0x65, 0x2c, 0xf7, 0xf5, 0x97, 0xd6, 0xf3,
	0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0x95, 0xf5, 0x0b, 0x18, 0x7e, 0x09, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TransactionClient is the client API for Transaction service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TransactionClient interface {
	ProxyOrderTranaction(ctx context.Context, in *ProxyOrderReq_DFB, opts ...grpc.CallOption) (*ProxyOrderResp_DFB, error)
	PayOrderTranaction(ctx context.Context, in *PayOrderRequest, opts ...grpc.CallOption) (*PayOrderResponse, error)
}

type transactionClient struct {
	cc *grpc.ClientConn
}

func NewTransactionClient(cc *grpc.ClientConn) TransactionClient {
	return &transactionClient{cc}
}

func (c *transactionClient) ProxyOrderTranaction(ctx context.Context, in *ProxyOrderReq_DFB, opts ...grpc.CallOption) (*ProxyOrderResp_DFB, error) {
	out := new(ProxyOrderResp_DFB)
	err := c.cc.Invoke(ctx, "/transaction.Transaction/ProxyOrderTranaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *transactionClient) PayOrderTranaction(ctx context.Context, in *PayOrderRequest, opts ...grpc.CallOption) (*PayOrderResponse, error) {
	out := new(PayOrderResponse)
	err := c.cc.Invoke(ctx, "/transaction.Transaction/PayOrderTranaction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TransactionServer is the server API for Transaction service.
type TransactionServer interface {
	ProxyOrderTranaction(context.Context, *ProxyOrderReq_DFB) (*ProxyOrderResp_DFB, error)
	PayOrderTranaction(context.Context, *PayOrderRequest) (*PayOrderResponse, error)
}

// UnimplementedTransactionServer can be embedded to have forward compatible implementations.
type UnimplementedTransactionServer struct {
}

func (*UnimplementedTransactionServer) ProxyOrderTranaction(ctx context.Context, req *ProxyOrderReq_DFB) (*ProxyOrderResp_DFB, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProxyOrderTranaction not implemented")
}
func (*UnimplementedTransactionServer) PayOrderTranaction(ctx context.Context, req *PayOrderRequest) (*PayOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PayOrderTranaction not implemented")
}

func RegisterTransactionServer(s *grpc.Server, srv TransactionServer) {
	s.RegisterService(&_Transaction_serviceDesc, srv)
}

func _Transaction_ProxyOrderTranaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProxyOrderReq_DFB)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).ProxyOrderTranaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transaction.Transaction/ProxyOrderTranaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).ProxyOrderTranaction(ctx, req.(*ProxyOrderReq_DFB))
	}
	return interceptor(ctx, in, info, handler)
}

func _Transaction_PayOrderTranaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PayOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TransactionServer).PayOrderTranaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/transaction.Transaction/PayOrderTranaction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TransactionServer).PayOrderTranaction(ctx, req.(*PayOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Transaction_serviceDesc = grpc.ServiceDesc{
	ServiceName: "transaction.Transaction",
	HandlerType: (*TransactionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ProxyOrderTranaction",
			Handler:    _Transaction_ProxyOrderTranaction_Handler,
		},
		{
			MethodName: "PayOrderTranaction",
			Handler:    _Transaction_PayOrderTranaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "transaction.proto",
}
