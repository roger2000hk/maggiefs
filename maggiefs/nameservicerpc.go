// Generated by go-rpcgen. Do not modify.

package maggiefs

import (
	"net/rpc"
)

type NameServiceService struct {
	impl NameService
}

func NewNameServiceService(impl NameService) *NameServiceService {
	return &NameServiceService{impl}
}

func RegisterNameServiceService(impl NameService) error {
	return rpc.RegisterName("NameService", NewNameServiceService(impl))
}

type NameServiceGetInodeRequest struct {
	Nodeid uint64
}

type NameServiceGetInodeResponse struct {
	I *Inode
}

func (s *NameServiceService) GetInode(request *NameServiceGetInodeRequest, response *NameServiceGetInodeResponse) (err error) {
	response.I, err = s.impl.GetInode(request.Nodeid)
	return
}

type NameServiceStatFsRequest struct {
}

type NameServiceStatFsResponse struct {
	Stat FsStat
}

func (s *NameServiceService) StatFs(request *NameServiceStatFsRequest, response *NameServiceStatFsResponse) (err error) {
	response.Stat, err = s.impl.StatFs()
	return
}

type NameServiceAddInodeRequest struct {
	Node Inode
}

type NameServiceAddInodeResponse struct {
	Id uint64
}

func (s *NameServiceService) AddInode(request *NameServiceAddInodeRequest, response *NameServiceAddInodeResponse) (err error) {
	response.Id, err = s.impl.AddInode(request.Node)
	return
}

type NameServiceSetInodeRequest struct {
	Node Inode
}

type NameServiceSetInodeResponse struct {
	NewNode *Inode
}

func (s *NameServiceService) SetInode(request *NameServiceSetInodeRequest, response *NameServiceSetInodeResponse) (err error) {
	response.NewNode, err = s.impl.SetInode(request.Node)
	return
}

type NameServiceLinkRequest struct {
	Parent uint64
	Child  uint64
	Name   string
	Force  bool
}

type NameServiceLinkResponse struct {
	NewNode *Inode
}

func (s *NameServiceService) Link(request *NameServiceLinkRequest, response *NameServiceLinkResponse) (err error) {
	response.NewNode, err = s.impl.Link(request.Parent, request.Child, request.Name, request.Force)
	return
}

type NameServiceUnlinkRequest struct {
	Parent uint64
	Name   string
}

type NameServiceUnlinkResponse struct {
	NewNode *Inode
}

func (s *NameServiceService) Unlink(request *NameServiceUnlinkRequest, response *NameServiceUnlinkResponse) (err error) {
	response.NewNode, err = s.impl.Unlink(request.Parent, request.Name)
	return
}

type NameServiceAddBlockRequest struct {
	Nodeid uint64
	Length uint32
}

type NameServiceAddBlockResponse struct {
	NewBlock Block
}

func (s *NameServiceService) AddBlock(request *NameServiceAddBlockRequest, response *NameServiceAddBlockResponse) (err error) {
	response.NewBlock, err = s.impl.AddBlock(request.Nodeid, request.Length)
	return
}

type NameServiceExtendBlockRequest struct {
	Nodeid  uint64
	BlockId uint64
	Delta   uint32
}

type NameServiceExtendBlockResponse struct {
	NewBlock Block
}

func (s *NameServiceService) ExtendBlock(request *NameServiceExtendBlockRequest, response *NameServiceExtendBlockResponse) (err error) {
	response.NewBlock, err = s.impl.ExtendBlock(request.Nodeid, request.BlockId, request.Delta)
	return
}

type NameServiceTruncateBlockRequest struct {
	BlockId uint64
	Delta   uint32
}

type NameServiceTruncateBlockResponse struct {
}

func (s *NameServiceService) TruncateBlock(request *NameServiceTruncateBlockRequest, response *NameServiceTruncateBlockResponse) (err error) {
	err = s.impl.TruncateBlock(request.BlockId, request.Delta)
	return
}

type NameServiceJoinRequest struct {
	DnId         int32
	NameDataAddr string
}

type NameServiceJoinResponse struct {
}

func (s *NameServiceService) Join(request *NameServiceJoinRequest, response *NameServiceJoinResponse) (err error) {
	err = s.impl.Join(request.DnId, request.NameDataAddr)
	return
}

type NameServiceNextVolIdRequest struct {
}

type NameServiceNextVolIdResponse struct {
	Id int32
}

func (s *NameServiceService) NextVolId(request *NameServiceNextVolIdRequest, response *NameServiceNextVolIdResponse) (err error) {
	response.Id, err = s.impl.NextVolId()
	return
}

type NameServiceNextDnIdRequest struct {
}

type NameServiceNextDnIdResponse struct {
	Id int32
}

func (s *NameServiceService) NextDnId(request *NameServiceNextDnIdRequest, response *NameServiceNextDnIdResponse) (err error) {
	response.Id, err = s.impl.NextDnId()
	return
}

type NameServiceRegisterVolRequest struct {
	DnId int32
	Stat VolumeStat
}

type NameServiceRegisterVolResponse struct {
}

func (s *NameServiceService) RegisterVol(request *NameServiceRegisterVolRequest, response *NameServiceRegisterVolResponse) (err error) {
	err = s.impl.RegisterVol(request.DnId, request.Stat)
	return
}

type NameServiceClient struct {
	client  *rpc.Client
	service string
}

func NewNameServiceClient(client *rpc.Client) *NameServiceClient {
	return &NameServiceClient{client, "NameService"}
}

func (_c *NameServiceClient) GetInode(nodeid uint64) (i *Inode, err error) {
	_request := &NameServiceGetInodeRequest{nodeid}
	_response := &NameServiceGetInodeResponse{}
	err = _c.client.Call(_c.service+".GetInode", _request, _response)
	return _response.I, err
}

func (_c *NameServiceClient) StatFs() (stat FsStat, err error) {
	_request := &NameServiceStatFsRequest{}
	_response := &NameServiceStatFsResponse{}
	err = _c.client.Call(_c.service+".StatFs", _request, _response)
	return _response.Stat, err
}

func (_c *NameServiceClient) AddInode(node Inode) (id uint64, err error) {
	_request := &NameServiceAddInodeRequest{node}
	_response := &NameServiceAddInodeResponse{}
	err = _c.client.Call(_c.service+".AddInode", _request, _response)
	return _response.Id, err
}

func (_c *NameServiceClient) SetInode(node Inode) (newNode *Inode, err error) {
	_request := &NameServiceSetInodeRequest{node}
	_response := &NameServiceSetInodeResponse{}
	err = _c.client.Call(_c.service+".SetInode", _request, _response)
	return _response.NewNode, err
}

func (_c *NameServiceClient) Link(parent uint64, child uint64, name string, force bool) (newNode *Inode, err error) {
	_request := &NameServiceLinkRequest{parent, child, name, force}
	_response := &NameServiceLinkResponse{}
	err = _c.client.Call(_c.service+".Link", _request, _response)
	return _response.NewNode, err
}

func (_c *NameServiceClient) Unlink(parent uint64, name string) (newNode *Inode, err error) {
	_request := &NameServiceUnlinkRequest{parent, name}
	_response := &NameServiceUnlinkResponse{}
	err = _c.client.Call(_c.service+".Unlink", _request, _response)
	return _response.NewNode, err
}

func (_c *NameServiceClient) AddBlock(nodeid uint64, length uint32) (newBlock Block, err error) {
	_request := &NameServiceAddBlockRequest{nodeid, length}
	_response := &NameServiceAddBlockResponse{}
	err = _c.client.Call(_c.service+".AddBlock", _request, _response)
	return _response.NewBlock, err
}

func (_c *NameServiceClient) ExtendBlock(nodeid uint64, blockId uint64, delta uint32) (newBlock Block, err error) {
	_request := &NameServiceExtendBlockRequest{nodeid, blockId, delta}
	_response := &NameServiceExtendBlockResponse{}
	err = _c.client.Call(_c.service+".ExtendBlock", _request, _response)
	return _response.NewBlock, err
}

func (_c *NameServiceClient) TruncateBlock(blockId uint64, delta uint32) (err error) {
	_request := &NameServiceTruncateBlockRequest{blockId, delta}
	_response := &NameServiceTruncateBlockResponse{}
	err = _c.client.Call(_c.service+".TruncateBlock", _request, _response)
	return err
}

func (_c *NameServiceClient) Join(dnId int32, nameDataAddr string) (err error) {
	_request := &NameServiceJoinRequest{dnId, nameDataAddr}
	_response := &NameServiceJoinResponse{}
	err = _c.client.Call(_c.service+".Join", _request, _response)
	return err
}

func (_c *NameServiceClient) NextVolId() (id int32, err error) {
	_request := &NameServiceNextVolIdRequest{}
	_response := &NameServiceNextVolIdResponse{}
	err = _c.client.Call(_c.service+".NextVolId", _request, _response)
	return _response.Id, err
}

func (_c *NameServiceClient) NextDnId() (id int32, err error) {
	_request := &NameServiceNextDnIdRequest{}
	_response := &NameServiceNextDnIdResponse{}
	err = _c.client.Call(_c.service+".NextDnId", _request, _response)
	return _response.Id, err
}

func (_c *NameServiceClient) RegisterVol(dnId int32, stat VolumeStat) (err error) {
	_request := &NameServiceRegisterVolRequest{dnId, stat}
	_response := &NameServiceRegisterVolResponse{}
	err = _c.client.Call(_c.service+".RegisterVol", _request, _response)
	return err
}
