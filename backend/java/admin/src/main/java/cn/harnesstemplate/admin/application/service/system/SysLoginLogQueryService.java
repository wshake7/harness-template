package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysLoginLog;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysLoginLogProxy;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class SysLoginLogQueryService {

    private final EasyQuerySupport easyQuery;

    public SysLoginLogQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<SysLoginLog> page(int pageNum, int pageSize, String username) {
        int offset = (pageNum - 1) * pageSize;
        return easyQuery.queryable(SysLoginLogProxy.createTable())
                .where(log -> log.username().like(username != null && !username.isEmpty(), username))
                .orderBy(log -> log.id().desc())
                .limit(offset, pageSize)
                .toList();
    }

    public long count(String username) {
        return easyQuery.queryable(SysLoginLogProxy.createTable())
                .where(log -> log.username().like(username != null && !username.isEmpty(), username))
                .count();
    }
}
