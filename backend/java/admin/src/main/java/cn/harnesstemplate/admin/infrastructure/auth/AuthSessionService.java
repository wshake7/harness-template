package cn.harnesstemplate.admin.infrastructure.auth;

import cn.dev33.satoken.stp.StpUtil;

/**
 * Wraps Sa-Token session operations, equivalent to Go's AuthService + Session.
 */
public class AuthSessionService {

    /**
     * Login and return a token, matching Go's Auth.Login(id).
     */
    public String login(Long userId) {
        StpUtil.login(userId);
        return StpUtil.getTokenValue();
    }

    /**
     * Logout by login ID, matching Go's Auth.Logout(loginID).
     */
    public void logout(Long loginId) {
        StpUtil.logout(loginId);
    }

    /**
     * Save session info, matching Go's session.SaveInfo(info).
     */
    public void saveSession(Long loginId, SessionInfo info) {
        StpUtil.getSessionByLoginId(loginId).set(SessionInfo.SESSION_KEY, info);
    }

    /**
     * Get session info from current session, matching Go's session.GetInfo().
     */
    public SessionInfo getSessionInfo() {
        return (SessionInfo) StpUtil.getSession().get(SessionInfo.SESSION_KEY);
    }

    /**
     * Get session info by login ID.
     */
    public SessionInfo getSessionInfo(Long loginId) {
        return (SessionInfo) StpUtil.getSessionByLoginId(loginId).get(SessionInfo.SESSION_KEY);
    }

    /**
     * Check if current request is authenticated.
     */
    public boolean isAuthenticated() {
        return StpUtil.isLogin();
    }

    /**
     * Get current login ID.
     */
    public Long getCurrentLoginId() {
        return StpUtil.getLoginIdAsLong();
    }
}
