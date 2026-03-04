package service

import (
	"context"
	"fmt"

	"github.com/Zara1024/AIOps/cloudops-server/internal/cmdb/model"
	"github.com/Zara1024/AIOps/cloudops-server/internal/cmdb/repository"
	"github.com/Zara1024/AIOps/cloudops-server/pkg/crypto"
	appErrors "github.com/Zara1024/AIOps/cloudops-server/pkg/errors"
)

// CMDBService CMDB 服务
type CMDBService struct {
	hostRepo     *repository.HostRepository
	groupRepo    *repository.HostGroupRepository
	sshRecordRepo *repository.SSHRecordRepository
}

// NewCMDBService 创建 CMDB 服务
func NewCMDBService(
	hostRepo *repository.HostRepository,
	groupRepo *repository.HostGroupRepository,
	sshRecordRepo *repository.SSHRecordRepository,
) *CMDBService {
	return &CMDBService{
		hostRepo:     hostRepo,
		groupRepo:    groupRepo,
		sshRecordRepo: sshRecordRepo,
	}
}

// ============= 主机管理 =============

// CreateHostRequest 创建主机请求
type CreateHostRequest struct {
	Hostname    string `json:"hostname" binding:"required,max=128"`
	IP          string `json:"ip" binding:"required"`
	Port        int    `json:"port"`
	OSType      string `json:"os_type"`
	AuthType    string `json:"auth_type"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PrivateKey  string `json:"private_key"`
	GroupID     *int64 `json:"group_id"`
	Description string `json:"description"`
	Labels      string `json:"labels"`
}

// CreateHost 创建主机
func (s *CMDBService) CreateHost(ctx context.Context, req *CreateHostRequest, createdBy int64) (*model.Host, error) {
	// 检查IP是否已存在
	exist, err := s.hostRepo.CheckIPExist(ctx, req.IP)
	if err != nil {
		return nil, appErrors.ErrInternal
	}
	if exist {
		return nil, appErrors.ErrDuplicate.WithMessage("IP地址已存在")
	}

	host := &model.Host{
		Hostname:    req.Hostname,
		IP:          req.IP,
		Port:        req.Port,
		OSType:      req.OSType,
		AuthType:    req.AuthType,
		Username:    req.Username,
		GroupID:      req.GroupID,
		Description: req.Description,
		Labels:      req.Labels,
		Status:      3, // 默认未知状态
		CreatedBy:   createdBy,
	}

	if host.Port == 0 {
		host.Port = 22
	}
	if host.Username == "" {
		host.Username = "root"
	}
	if host.OSType == "" {
		host.OSType = "linux"
	}
	if host.AuthType == "" {
		host.AuthType = "password"
	}

	// 加密存储密码/密钥
	if req.Password != "" {
		encrypted, err := crypto.AESEncrypt(req.Password)
		if err != nil {
			return nil, appErrors.ErrInternal.WithDetail("密码加密失败")
		}
		host.Password = encrypted
	}
	if req.PrivateKey != "" {
		encrypted, err := crypto.AESEncrypt(req.PrivateKey)
		if err != nil {
			return nil, appErrors.ErrInternal.WithDetail("密钥加密失败")
		}
		host.PrivateKey = encrypted
	}

	if err := s.hostRepo.Create(ctx, host); err != nil {
		return nil, appErrors.ErrInternal.WithDetail(err.Error())
	}

	return host, nil
}

// UpdateHostRequest 更新主机请求
type UpdateHostRequest struct {
	Hostname    string `json:"hostname"`
	Port        int    `json:"port"`
	OSType      string `json:"os_type"`
	AuthType    string `json:"auth_type"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	PrivateKey  string `json:"private_key"`
	GroupID     *int64 `json:"group_id"`
	Description string `json:"description"`
	Labels      string `json:"labels"`
	Status      int    `json:"status"`
}

// UpdateHost 更新主机
func (s *CMDBService) UpdateHost(ctx context.Context, id int64, req *UpdateHostRequest) error {
	host, err := s.hostRepo.FindByID(ctx, id)
	if err != nil {
		return appErrors.ErrHostNotFound
	}

	if req.Hostname != "" {
		host.Hostname = req.Hostname
	}
	if req.Port > 0 {
		host.Port = req.Port
	}
	if req.OSType != "" {
		host.OSType = req.OSType
	}
	if req.AuthType != "" {
		host.AuthType = req.AuthType
	}
	if req.Username != "" {
		host.Username = req.Username
	}
	if req.Password != "" {
		encrypted, _ := crypto.AESEncrypt(req.Password)
		host.Password = encrypted
	}
	if req.PrivateKey != "" {
		encrypted, _ := crypto.AESEncrypt(req.PrivateKey)
		host.PrivateKey = encrypted
	}
	host.GroupID = req.GroupID
	if req.Description != "" {
		host.Description = req.Description
	}
	if req.Labels != "" {
		host.Labels = req.Labels
	}
	if req.Status > 0 {
		host.Status = req.Status
	}

	return s.hostRepo.Update(ctx, host)
}

// DeleteHost 删除主机
func (s *CMDBService) DeleteHost(ctx context.Context, id int64) error {
	return s.hostRepo.Delete(ctx, id)
}

// BatchDeleteHosts 批量删除
func (s *CMDBService) BatchDeleteHosts(ctx context.Context, ids []int64) error {
	return s.hostRepo.BatchDelete(ctx, ids)
}

// GetHost 获取主机详情
func (s *CMDBService) GetHost(ctx context.Context, id int64) (*model.Host, error) {
	return s.hostRepo.FindByID(ctx, id)
}

// ListHosts 主机列表
func (s *CMDBService) ListHosts(ctx context.Context, page, pageSize int, ip, hostname string, status, groupID int) ([]model.Host, int64, error) {
	return s.hostRepo.List(ctx, page, pageSize, ip, hostname, status, groupID)
}

// BatchUpdateGroup 批量更新分组
func (s *CMDBService) BatchUpdateGroup(ctx context.Context, ids []int64, groupID int64) error {
	return s.hostRepo.BatchUpdateGroup(ctx, ids, groupID)
}

// GetHostSSHCredentials 获取主机SSH凭据（解密）
func (s *CMDBService) GetHostSSHCredentials(ctx context.Context, id int64) (username, password, privateKey string, port int, ip string, err error) {
	host, err := s.hostRepo.FindByID(ctx, id)
	if err != nil {
		return "", "", "", 0, "", appErrors.ErrHostNotFound
	}

	username = host.Username
	port = host.Port
	ip = host.IP

	if host.Password != "" {
		password, err = crypto.AESDecrypt(host.Password)
		if err != nil {
			return "", "", "", 0, "", appErrors.ErrInternal.WithDetail("解密密码失败")
		}
	}
	if host.PrivateKey != "" {
		privateKey, err = crypto.AESDecrypt(host.PrivateKey)
		if err != nil {
			return "", "", "", 0, "", appErrors.ErrInternal.WithDetail("解密密钥失败")
		}
	}

	return username, password, privateKey, port, ip, nil
}

// ============= 主机分组 =============

// GetGroupTree 获取分组树
func (s *CMDBService) GetGroupTree(ctx context.Context) ([]*model.HostGroup, error) {
	groups, err := s.groupRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 为每个分组统计主机数
	for i := range groups {
		count, _ := s.groupRepo.CountHostsByGroupID(ctx, groups[i].ID)
		groups[i].HostCount = count
	}

	return buildGroupTree(groups, 0), nil
}

// CreateGroup 创建分组
func (s *CMDBService) CreateGroup(ctx context.Context, group *model.HostGroup) error {
	return s.groupRepo.Create(ctx, group)
}

// UpdateGroup 更新分组
func (s *CMDBService) UpdateGroup(ctx context.Context, group *model.HostGroup) error {
	return s.groupRepo.Update(ctx, group)
}

// DeleteGroup 删除分组
func (s *CMDBService) DeleteGroup(ctx context.Context, id int64) error {
	count, err := s.groupRepo.CountHostsByGroupID(ctx, id)
	if err != nil {
		return appErrors.ErrInternal
	}
	if count > 0 {
		return appErrors.ErrBadRequest.WithMessage(fmt.Sprintf("分组下还有 %d 台主机，不能删除", count))
	}
	return s.groupRepo.Delete(ctx, id)
}

// ============= SSH 记录 =============

// ListSSHRecords 查询SSH操作记录
func (s *CMDBService) ListSSHRecords(ctx context.Context, page, pageSize int, hostID int64, username string) ([]model.SSHRecord, int64, error) {
	return s.sshRecordRepo.List(ctx, page, pageSize, hostID, username)
}

// 树形构建
func buildGroupTree(groups []model.HostGroup, parentID int64) []*model.HostGroup {
	var tree []*model.HostGroup
	for i := range groups {
		if groups[i].ParentID == parentID {
			node := &groups[i]
			node.Children = buildGroupTree(groups, node.ID)
			tree = append(tree, node)
		}
	}
	return tree
}
