package cn.harnesstemplate.admin.infrastructure.persistence;

import com.easy.query.api.proxy.client.DefaultEasyEntityQuery;
import com.easy.query.api.proxy.client.EasyEntityQuery;
import com.easy.query.api.proxy.entity.select.EntityQueryable;
import com.easy.query.api.proxy.entity.update.ExpressionUpdatable;
import com.easy.query.core.api.client.EasyQueryClient;
import com.easy.query.core.basic.api.select.ClientQueryable;
import com.easy.query.core.basic.api.update.ClientExpressionUpdatable;
import com.easy.query.core.proxy.ProxyEntity;
import org.springframework.stereotype.Component;

import java.time.LocalDateTime;
import java.util.Arrays;
import java.util.List;

@Component
public class EasyQuerySupport {

    private final EasyQueryClient client;
    private final EasyEntityQuery entityQuery;

    public EasyQuerySupport(EasyQueryClient client) {
        this.client = client;
        this.entityQuery = new DefaultEasyEntityQuery(client);
    }

    public <T> List<T> query(String sql, Class<T> resultType, Object... args) {
        return client.sqlQuery(sql, resultType, Arrays.asList(args));
    }

    public <T> ClientQueryable<T> queryable(Class<T> entityType) {
        return client.queryable(entityType);
    }

    public <TProxy extends ProxyEntity<TProxy, T>, T> EntityQueryable<TProxy, T> queryable(TProxy proxy) {
        return entityQuery.queryable(proxy);
    }

    public <T> List<T> queryForList(String sql, Class<T> resultType, Object... args) {
        return query(sql, resultType, args);
    }

    public <T> T queryForObject(String sql, Class<T> resultType, Object... args) {
        List<T> rows = query(sql, resultType, args);
        return rows.isEmpty() ? null : rows.get(0);
    }

    public long execute(String sql, Object... args) {
        return client.sqlExecute(sql, Arrays.asList(args));
    }

    public <T> T insert(T entity) {
        client.insertable(entity).executeRows(true);
        return entity;
    }

    public <T> ClientExpressionUpdatable<T> updatable(Class<T> entityType) {
        return client.updatable(entityType);
    }

    public <TProxy extends ProxyEntity<TProxy, T>, T> ExpressionUpdatable<TProxy, T> expressionUpdatable(TProxy proxy) {
        return entityQuery.expressionUpdatable(proxy);
    }

    public <T> long update(T entity) {
        return client.updatable(entity).executeRows();
    }

    public long softDelete(Class<?> entityType, Long id) {
        return client.updatable(entityType)
                .set("deletedAt", System.currentTimeMillis())
                .where(w -> {
                    w.eq("id", id);
                    w.eq("deletedAt", 0L);
                })
                .executeRows();
    }

    public long markUpdatedAt(Class<?> entityType, Long id) {
        return client.updatable(entityType)
                .set("updatedAt", LocalDateTime.now())
                .whereById(id)
                .executeRows();
    }
}
