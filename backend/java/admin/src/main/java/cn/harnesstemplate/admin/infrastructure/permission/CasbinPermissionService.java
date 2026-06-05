package cn.harnesstemplate.admin.infrastructure.permission;

import org.casbin.jcasbin.main.SyncedEnforcer;
import org.casbin.jcasbin.model.Model;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.jdbc.core.JdbcTemplate;

import javax.sql.DataSource;
import java.util.ArrayList;
import java.util.List;

/**
 * Casbin permission service matching Go's casbin package behavior.
 * Loads model from sys_casbin_model table, uses JDBC adapter for policy
 * persistence, and provides enforce with OR logic across user/role subjects.
 */
public class CasbinPermissionService {

    private static final Logger log = LoggerFactory.getLogger(CasbinPermissionService.class);

    /**
     * Fallback model matching Go's normalizeModelContent output.
     * Used when sys_casbin_model table has no enabled model.
     */
    static final String DEFAULT_MODEL_TEXT = """
            [request_definition]
            r = sub, obj, act

            [policy_definition]
            p = sub, obj, act

            [policy_effect]
            e = some(where (p.eft == allow))

            [matchers]
            m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && r.act == p.act
            """;

    private final SyncedEnforcer enforcer;

    /**
     * Create service with database-backed enforcer.
     * Loads model from sys_casbin_model table and policies from casbin_rule table.
     * Falls back to in-memory mode if database tables are not available.
     */
    public CasbinPermissionService(DataSource dataSource) {
        Model model = new Model();
        String modelText = loadModelText(dataSource);
        model.loadModelFromText(modelText);

        SyncedEnforcer enforcer;
        try {
            CasbinJdbcAdapter adapter = new CasbinJdbcAdapter(dataSource);
            adapter.loadPolicy(model);
            enforcer = new SyncedEnforcer(model, adapter);
            enforcer.enableAutoSave(true);
            normalizeLegacyPolicies(enforcer);
            log.info("Casbin enforcer initialized (db-backed) with {} policies", enforcer.getPolicy().size());
        } catch (Exception e) {
            log.warn("Falling back to in-memory casbin enforcer: {}", e.getMessage());
            enforcer = new SyncedEnforcer(model);
        }
        this.enforcer = enforcer;
    }

    /**
     * Test-only constructor: creates in-memory enforcer without database.
     */
    public CasbinPermissionService() {
        Model model = new Model();
        model.loadModelFromText(DEFAULT_MODEL_TEXT);
        this.enforcer = new SyncedEnforcer(model);
    }

    /**
     * Test-only constructor with pre-configured enforcer.
     */
    CasbinPermissionService(SyncedEnforcer enforcer) {
        this.enforcer = enforcer;
    }

    // ---- Policy management ----

    /**
     * Add a policy rule for a subject (role:code or user:id) with path and method.
     */
    public void addPolicy(String sub, String obj, String act) {
        enforcer.addPolicy(sub, obj, act);
    }

    /**
     * Remove a policy rule.
     */
    public void removePolicy(String sub, String obj, String act) {
        enforcer.removePolicy(sub, obj, act);
    }

    /**
     * Remove all policies for a given subject.
     */
    public void removePoliciesForSubject(String sub) {
        List<List<String>> toRemove = new ArrayList<>();
        for (List<String> policy : enforcer.getPolicy()) {
            if (policy.size() >= 3 && policy.get(0).equals(sub)) {
                toRemove.add(policy);
            }
        }
        for (List<String> p : toRemove) {
            enforcer.removePolicy(p);
        }
    }

    /**
     * Batch add policies matching Go's AddPoliciesEx.
     */
    public void addPolicies(List<List<String>> rules) {
        if (rules == null || rules.isEmpty()) return;
        enforcer.addPolicies(rules);
    }

    /**
     * Remove policies matching Go's RemoveFilteredPolicy.
     */
    public void removeFilteredPolicy(int fieldIndex, String... fieldValues) {
        enforcer.removeFilteredPolicy(fieldIndex, fieldValues);
    }

    // ---- Enforcement ----

    /**
     * Build subjects matching Go's pattern: user:id + role:code for each role.
     */
    public static List<String> buildSubjects(long userId, List<String> roleCodes) {
        List<String> subjects = new ArrayList<>();
        subjects.add("user:" + userId);
        if (roleCodes != null) {
            for (String role : roleCodes) {
                String trimmed = role.trim();
                if (!trimmed.isEmpty()) {
                    subjects.add("role:" + trimmed);
                }
            }
        }
        return subjects;
    }

    /**
     * Enforce permission check with OR logic (any subject pass = allowed).
     * Matches Go's CasbinAPIMiddleware behavior.
     */
    public boolean enforce(List<String> subjects, String path, String method) {
        String act = method.toUpperCase();
        for (String sub : subjects) {
            boolean ok;
            try {
                ok = enforcer.enforce(sub, path, act);
            } catch (Exception e) {
                log.debug("Casbin enforce error for sub={} obj={} act={}: {}",
                        sub, path, act, e.getMessage());
                continue;
            }
            if (ok) {
                return true;
            }
        }
        log.debug("Casbin auth denied: subjects={} path={} method={}",
                subjects, path, act);
        return false;
    }

    public SyncedEnforcer getEnforcer() {
        return enforcer;
    }

    // ---- Initialization helpers matching Go's casbin.go ----

    private static String loadModelText(DataSource dataSource) {
        JdbcTemplate jdbc = new JdbcTemplate(dataSource);
        try {
            List<String> contents = jdbc.query(
                    "SELECT content FROM sys_casbin_model WHERE is_enabled = true LIMIT 1",
                    (rs, rowNum) -> rs.getString("content"));
            if (!contents.isEmpty()) {
                return normalizeModelContent(contents.get(0));
            }
        } catch (Exception e) {
            log.warn("Could not load casbin model from sys_casbin_model: {}. Using default model.",
                    e.getMessage());
        }
        log.info("No enabled casbin model found in sys_casbin_model, using default");
        return DEFAULT_MODEL_TEXT;
    }

    /**
     * Normalize legacy model content matching Go's normalizeModelContent.
     */
    static String normalizeModelContent(String content) {
        content = content.replace("p = sub_rule, obj_rule, act", "p = sub, obj, act");
        content = content.replace(
                "m = eval(p.sub_rule) && eval(p.obj_rule) && r.act == p.act",
                "m = r.sub == p.sub && keyMatch2(r.obj, p.obj) && r.act == p.act");
        return content;
    }

    /**
     * Normalize legacy policies matching Go's normalizeLegacyPolicies.
     */
    static void normalizeLegacyPolicies(SyncedEnforcer enforcer) {
        List<List<String>> policies;
        try {
            policies = enforcer.getPolicy();
        } catch (Exception e) {
            log.warn("Failed to get policies for normalization", e);
            return;
        }
        for (List<String> policy : policies) {
            if (policy.size() < 3) continue;
            String sub = policy.get(0);
            String obj = policy.get(1);
            String newSub = parseLegacySubjectRule(sub);
            String newObj = parseLegacyObjectRule(obj);
            if (newSub == null && newObj == null) continue;

            List<String> next = new ArrayList<>(policy);
            if (newSub != null) next.set(0, newSub);
            if (newObj != null) next.set(1, newObj);
            enforcer.removePolicy(policy);
            enforcer.addPolicy(next);
        }
    }

    static String parseLegacySubjectRule(String value) {
        String trimmed = value.trim();
        String prefix = "r.sub == ";
        if (!trimmed.startsWith(prefix)) return null;
        String quoted = trimmed.substring(prefix.length()).trim();
        return unquote(quoted);
    }

    static String parseLegacyObjectRule(String value) {
        String trimmed = value.trim();
        String prefix = "keyMatch2(r.obj, ";
        if (!trimmed.startsWith(prefix) || !trimmed.endsWith(")")) return null;
        String quoted = trimmed.substring(prefix.length(), trimmed.length() - 1).trim();
        return unquote(quoted);
    }

    private static String unquote(String s) {
        if (s.length() >= 2 && s.startsWith("\"") && s.endsWith("\"")) {
            return s.substring(1, s.length() - 1);
        }
        return null;
    }
}
