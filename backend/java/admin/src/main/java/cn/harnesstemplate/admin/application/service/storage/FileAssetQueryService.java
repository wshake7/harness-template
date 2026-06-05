package cn.harnesstemplate.admin.application.service.storage;

import cn.harnesstemplate.admin.domain.entity.FileAsset;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.FileAssetProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class FileAssetQueryService {

    private final EasyQuerySupport easyQuery;

    public FileAssetQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<FileAsset> page(int pageNum, int pageSize) {
        int offset = (pageNum - 1) * pageSize;
        return easyQuery.queryable(FileAssetProxy.createTable())
                .where(file -> file.deletedAt().eq(0L))
                .orderBy(file -> file.id().desc())
                .limit(offset, pageSize)
                .toList();
    }

    public long count() {
        return easyQuery.queryable(FileAssetProxy.createTable())
                .where(file -> file.deletedAt().eq(0L))
                .count();
    }

    public FileAsset findById(Long id) {
        return easyQuery.queryable(FileAssetProxy.createTable())
                .where(file -> {
                    file.id().eq(id);
                    file.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public FileAsset save(FileAsset f) {
        return easyQuery.insert(f);
    }

    public void updateStatus(Long id, String status) {
        easyQuery.expressionUpdatable(FileAssetProxy.createTable())
                .setColumns(file -> {
                    file.status().set(status);
                    file.updatedAt().set(LocalDateTime.now());
                })
                .where(file -> file.id().eq(id))
                .executeRows();
    }

    public void delete(Long id) {
        easyQuery.expressionUpdatable(FileAssetProxy.createTable())
                .setColumns(file -> {
                    file.status().set("deleted");
                    file.deletedAt().set(System.currentTimeMillis());
                    file.updatedAt().set(LocalDateTime.now());
                })
                .where(file -> file.id().eq(id))
                .executeRows();
    }
}
