package cn.harnesstemplate.admin.infrastructure.auth;

import java.util.List;

/**
 * Session info stored in Sa-Token session, equivalent to Go's auth.SessionInfo.
 */
public class SessionInfo {

    public static final String SESSION_KEY = "sessionInfo";

    private String privateKey;
    private Long id;
    private String username;
    private List<String> roleCodes;
    private List<Long> roleIds;
    private String language;

    public SessionInfo() {}

    public SessionInfo(String privateKey, Long id, String username,
                       List<String> roleCodes, List<Long> roleIds) {
        this.privateKey = privateKey;
        this.id = id;
        this.username = username;
        this.roleCodes = roleCodes;
        this.roleIds = roleIds;
    }

    public String getPrivateKey() { return privateKey; }
    public void setPrivateKey(String privateKey) { this.privateKey = privateKey; }
    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public String getUsername() { return username; }
    public void setUsername(String username) { this.username = username; }
    public List<String> getRoleCodes() { return roleCodes; }
    public void setRoleCodes(List<String> roleCodes) { this.roleCodes = roleCodes; }
    public List<Long> getRoleIds() { return roleIds; }
    public void setRoleIds(List<Long> roleIds) { this.roleIds = roleIds; }
    public String getLanguage() { return language; }
    public void setLanguage(String language) { this.language = language; }
}
