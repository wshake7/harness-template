package logic

import (
	"admin/internal/auth"
	"admin/internal/domains"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/service"
	"admin/internal/services/orm/query"
	"errors"

	"go-common/utils/encrypt/rsa_util"
	"go-common/utils/passwd"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AccountHandler struct {
	Q           *query.Query
	Auth        service.AuthService
	LoginLogger service.LoginLogger
}

func NewAccountHandler(q *query.Query, auth service.AuthService, logger service.LoginLogger) *AccountHandler {
	return &AccountHandler{Q: q, Auth: auth, LoginLogger: logger}
}

type ReqAccountPwdLogin struct {
	Username string `json:"username" binding:"required,max=24" binding_msg:"required=用户名不能为空,max=用户名最多24位"`
	Pwd      string `json:"pwd" binding:"required,min=6" binding_msg:"required=密码不能为空,min=密码最少6位"`
}

type ResAccountPwdLogin struct {
	Token     string `json:"token"`
	PublicKey string `json:"publicKey"`
}

// @Summary 用户密码登录
// @Remark 使用用户名和密码登录，返回访问令牌与本次会话使用的加密公钥
// @Tags Account
// @Accept json
// @Produce json
// @Param req body ReqAccountPwdLogin true "登录参数"
// @Success 200 {object} res.Response{data=ResAccountPwdLogin} "成功"
// @Router /api/account/login/pwd [post]
func (h *AccountHandler) PwdLogin(ctx *handler.Ctx, req *ReqAccountPwdLogin) (*ResAccountPwdLogin, error) {
	logger := ctx.L().With(zap.String("username", req.Username))

	sysUser := h.Q.SysUser
	result, err := sysUser.Where(sysUser.Username.Eq(req.Username)).Select(sysUser.ID, sysUser.Username, sysUser.Password).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.LoginLogger.RecordPwdLogin(ctx, req.Username, nil, domains.StatusLoginFail, false, "invalid username or password")
			return nil, errors.New("用户名或密码无效")
		}
		logger.Error("获取用户失败", zap.Error(err))
		h.LoginLogger.RecordPwdLogin(ctx, req.Username, nil, domains.StatusFail, false, "query user failed")
		return nil, errors.New("登录失败")
	}

	if !passwd.Match(req.Pwd, result.Password) {
		userID := result.ID
		h.LoginLogger.RecordPwdLogin(ctx, req.Username, &userID, domains.StatusLoginFail, false, "invalid username or password")
		return nil, errors.New("用户名或密码无效")
	}

	sysRole := h.Q.SysRole
	sysUserRole := h.Q.SysUserRole
	userRoles, err := sysUserRole.Select(sysUserRole.ID, sysUserRole.UserID, sysUserRole.RoleID).Preload(sysUserRole.SysRole.Select(sysRole.ID, sysRole.Code).On(sysRole.IsEnabled.Is(true))).Where(sysUserRole.UserID.Eq(result.ID)).Find()
	if err != nil {
		logger.Error("获取用户角色失败", zap.Error(err), zap.Uint64("userID", result.ID))
		userID := result.ID
		h.LoginLogger.RecordPwdLogin(ctx, result.Username, &userID, domains.StatusFail, false, "query user roles failed")
		return nil, errors.New("登录失败")
	}

	roleCodes := make([]string, 0, len(userRoles))
	roleIDs := make([]uint64, 0, len(userRoles))
	for _, userRole := range userRoles {
		if userRole == nil || userRole.SysRole == nil {
			continue
		}
		role := userRole.SysRole
		if role.Code != "" {
			roleCodes = append(roleCodes, role.Code)
		}
		if role.ID > 0 {
			roleIDs = append(roleIDs, role.ID)
		}
	}

	token, err := h.Auth.Login(result.ID)
	if err != nil {
		logger.Error("获取token失败", zap.Error(err))
		userID := result.ID
		h.LoginLogger.RecordPwdLogin(ctx, result.Username, &userID, domains.StatusLoginFail, false, "create token failed")
		return nil, errors.New("登录失败")
	}

	privateKey, publicKey, err := rsa_util.GenerateKeyPair()
	if err != nil {
		logger.Error("获取rsaKey错误", zap.Error(err))
		userID := result.ID
		h.LoginLogger.RecordPwdLogin(ctx, result.Username, &userID, domains.StatusFail, false, "generate rsa key failed")
		return nil, errors.New("登录失败")
	}

	err = h.Auth.SaveSession(result.ID, &auth.SessionInfo{
		PrivateKey: privateKey,
		Id:         result.ID,
		Username:   result.Username,
		RoleCodes:  roleCodes,
		RoleIDs:    roleIDs,
	})
	if err != nil {
		logger.Error("保存SessionInfo错误", zap.Error(err))
		userID := result.ID
		h.LoginLogger.RecordPwdLogin(ctx, result.Username, &userID, domains.StatusFail, false, "save session failed")
		return nil, errors.New("登录失败")
	}

	userID := result.ID
	h.LoginLogger.RecordPwdLogin(ctx, result.Username, &userID, domains.StatusOk, true, "")
	return &ResAccountPwdLogin{
		Token:     token,
		PublicKey: publicKey,
	}, nil
}

type ReqAccountLogout struct {
	Token string `cookie:"token" binding:"required" binding_msg:"required=请求错误"`
}

// @Summary 退出登录
// @Remark 根据当前登录 token 退出会话
// @Tags Account
// @Produce json
// @Param Authorization header string false "登录 Token，或使用 cookie 中的 token"
// @Success 200 {object} res.Response "成功"
// @Router /api/account/logout [get]
func (h *AccountHandler) Logout(ctx *handler.Ctx, req *ReqAccountLogout) error {
	loginID, err := h.Auth.GetLoginID(req.Token)
	if err != nil {
		ctx.L().Error("获取loginId失败", zap.Error(err))
		return auth.CheckLoginErr(err)
	}
	err = h.Auth.Logout(loginID)
	if err != nil {
		ctx.L().Error("退出登录失败", zap.Error(err), zap.String("token", req.Token))
		return auth.CheckLoginErr(err)
	}
	return nil
}

type ReqAccountChangePwd struct {
	OldPwd string `json:"oldPwd" binding:"required,min=6" binding_msg:"required=原始密码不能为空,min=原始密码最少6位"`
	NewPwd string `json:"newPwd" binding:"required,min=6" binding_msg:"required=新密码不能为空,min=新密码最少6位"`
}

// @Summary 修改密码
// @Remark 校验旧密码后更新当前登录用户密码
// @Tags Account
// @Accept json
// @Produce json
// @Param req body ReqAccountChangePwd true "修改密码参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/account/changePwd [post]
func (h *AccountHandler) ChangePwd(ctx *handler.Ctx, req *ReqAccountChangePwd) error {
	sysUser := h.Q.SysUser
	result, err := sysUser.Where(sysUser.ID.Eq(ctx.SessionInfo.Id)).Select(sysUser.ID, sysUser.Password).First()
	if err != nil {
		ctx.L().Error("获取用户密码失败", zap.Error(err))
		return res.FailDefault
	}

	if !passwd.Match(req.OldPwd, result.Password) {
		return errors.New("原密码错误")
	}

	encodePwd, err := passwd.Encode(req.NewPwd)
	if err != nil {
		ctx.L().Error("密码加密失败", zap.Error(err))
		return res.FailDefault
	}
	_, err = sysUser.Where(sysUser.ID.Eq(ctx.SessionInfo.Id)).Update(sysUser.Password, encodePwd)
	if err != nil {
		ctx.L().Error("修改密码失败", zap.Error(err))
		return res.FailDefault
	}
	return nil
}
