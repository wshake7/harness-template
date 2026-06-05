package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysApiLog;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysApiLogProxy;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class SysApiLogQueryService {

    private final EasyQuerySupport easyQuery;

    public SysApiLogQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<SysApiLog> page(int pageNum, int pageSize, String path, String method) {
        int offset = (pageNum - 1) * pageSize;
        return easyQuery.queryable(SysApiLogProxy.createTable())
                .where(log -> {
                    log.path().like(path != null && !path.isEmpty(), path);
                    log.method().eq(method != null && !method.isEmpty(), method == null ? null : method.toUpperCase());
                })
                .orderBy(log -> log.id().desc())
                .limit(offset, pageSize)
                .toList();
    }

    public long count(String path, String method) {
        return easyQuery.queryable(SysApiLogProxy.createTable())
                .where(log -> {
                    log.path().like(path != null && !path.isEmpty(), path);
                    log.method().eq(method != null && !method.isEmpty(), method == null ? null : method.toUpperCase());
                })
                .count();
    }
}
