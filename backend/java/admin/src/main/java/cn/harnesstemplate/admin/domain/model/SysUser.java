package cn.harnesstemplate.admin.domain.model;

import java.util.List;

/**
 * Minimal user model for authentication. Will be replaced by full entity in Milestone 3.
 */
public class SysUser {

    private Long id;
    private String username;
    private String password;
    private List<SysRole> roles;

    public SysUser() {}

    public SysUser(Long id, String username, String password) {
        this.id = id;
        this.username = username;
        this.password = password;
    }

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public String getUsername() { return username; }
    public void setUsername(String username) { this.username = username; }
    public String getPassword() { return password; }
    public void setPassword(String password) { this.password = password; }
    public List<SysRole> getRoles() { return roles; }
    public void setRoles(List<SysRole> roles) { this.roles = roles; }
}
