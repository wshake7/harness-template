package cn.harnesstemplate.admin.domain.model;

/**
 * Minimal role model for authentication. Will be replaced by full entity in Milestone 3.
 */
public class SysRole {

    private Long id;
    private String code;

    public SysRole() {}

    public SysRole(Long id, String code) {
        this.id = id;
        this.code = code;
    }

    public Long getId() { return id; }
    public void setId(Long id) { this.id = id; }
    public String getCode() { return code; }
    public void setCode(String code) { this.code = code; }
}
