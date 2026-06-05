package cn.harnesstemplate.admin.infrastructure.permission;

import org.casbin.jcasbin.model.Model;
import org.casbin.jcasbin.persist.Adapter;
import org.springframework.jdbc.core.JdbcTemplate;

import javax.sql.DataSource;
import java.util.*;
import java.util.stream.Collectors;

/**
 * JDBC adapter for jCasbin that reads/writes policies from the casbin_rule table.
 * Matches Go's gormadapter behavior.
 */
public class CasbinJdbcAdapter implements Adapter {

    private static final String TABLE_NAME = "casbin_rule";

    private final JdbcTemplate jdbc;

    public CasbinJdbcAdapter(DataSource dataSource) {
        this.jdbc = new JdbcTemplate(dataSource);
    }

    // ---- Load ----

    @Override
    public void loadPolicy(Model model) {
        List<RuleRow> rows = jdbc.query(
                "SELECT ptype, v0, v1, v2, v3, v4, v5 FROM " + TABLE_NAME,
                (rs, rowNum) -> new RuleRow(
                        rs.getString("ptype"),
                        rs.getString("v0"),
                        rs.getString("v1"),
                        rs.getString("v2"),
                        rs.getString("v3"),
                        rs.getString("v4"),
                        rs.getString("v5")
                ));
        for (RuleRow row : rows) {
            List<String> values = row.values();
            if (values.isEmpty()) continue;
            model.addPolicy(row.sec(), row.ptype, values);
        }
    }

    // ---- Save (full sync) ----

    @Override
    public void savePolicy(Model model) {
        jdbc.update("DELETE FROM " + TABLE_NAME);
        List<Object[]> batch = new ArrayList<>();
        for (Map.Entry<String, Map<String, org.casbin.jcasbin.model.Assertion>> secEntry :
                model.model.entrySet()) {
            String sec = secEntry.getKey();
            for (Map.Entry<String, org.casbin.jcasbin.model.Assertion> ptypeEntry :
                    secEntry.getValue().entrySet()) {
                String ptype = ptypeEntry.getKey();
                for (List<String> rule : ptypeEntry.getValue().policy) {
                    batch.add(toRow(ptype, rule));
                }
            }
        }
        jdbc.batchUpdate(
                "INSERT INTO " + TABLE_NAME + " (ptype, v0, v1, v2, v3, v4, v5) VALUES (?, ?, ?, ?, ?, ?, ?)",
                batch);
    }

    // ---- Single policy add / remove (used with auto-save) ----

    @Override
    public void addPolicy(String sec, String ptype, List<String> rule) {
        jdbc.update(
                "INSERT INTO " + TABLE_NAME + " (ptype, v0, v1, v2, v3, v4, v5) VALUES (?, ?, ?, ?, ?, ?, ?)",
                toRow(ptype, rule));
    }

    @Override
    public void removePolicy(String sec, String ptype, List<String> rule) {
        StringBuilder sql = new StringBuilder("DELETE FROM " + TABLE_NAME + " WHERE ptype = ?");
        Object[] values = rule.toArray(new String[0]);
        for (int i = 0; i < 6; i++) {
            if (i < values.length) {
                sql.append(" AND v").append(i).append(" = ?");
            } else {
                sql.append(" AND v").append(i).append(" = ''");
            }
        }
        Object[] params = new Object[7];
        params[0] = ptype;
        for (int i = 0; i < 6; i++) {
            if (i < values.length) {
                params[i + 1] = values[i];
            } else {
                params[i + 1] = "";
            }
        }
        jdbc.update(sql.toString(), params);
    }

    @Override
    public void removeFilteredPolicy(String sec, String ptype, int fieldIndex, String... fieldValues) {
        StringBuilder sql = new StringBuilder("DELETE FROM " + TABLE_NAME + " WHERE ptype = ?");
        List<String> params = new ArrayList<>();
        params.add(ptype);
        for (int i = 0; i < fieldValues.length; i++) {
            sql.append(" AND v").append(fieldIndex + i).append(" = ?");
            params.add(fieldValues[i]);
        }
        jdbc.update(sql.toString(), params.toArray());
    }

    // ---- Helpers ----

    private static Object[] toRow(String ptype, List<String> rule) {
        String[] row = new String[7];
        row[0] = ptype;
        for (int i = 0; i < 6; i++) {
            row[i + 1] = i < rule.size() ? rule.get(i) : "";
        }
        return row;
    }

    private record RuleRow(String ptype, String v0, String v1, String v2, String v3, String v4, String v5) {
        String sec() {
            return ptype != null ? ptype : "p";
        }

        List<String> values() {
            List<String> vals = new ArrayList<>();
            for (String v : new String[]{v0, v1, v2, v3, v4, v5}) {
                if (v != null && !v.isEmpty()) {
                    vals.add(v);
                } else {
                    break;
                }
            }
            return vals;
        }
    }
}
