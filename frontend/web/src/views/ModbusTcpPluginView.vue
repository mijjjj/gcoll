<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { NButton, NDescriptions, NDescriptionsItem, NTag } from 'naive-ui'
import { ArrowLeft, BookOpen, ListChecks, RefreshCw } from '@lucide/vue'
import { useRouter } from 'vue-router'
import PageHeader from '../components/common/PageHeader.vue'
import SectionPanel from '../components/common/SectionPanel.vue'
import StatusBadge from '../components/common/StatusBadge.vue'
import { useConsoleStore } from '../stores/console'

const router = useRouter()
const consoleStore = useConsoleStore()

const plugin = computed(() => consoleStore.plugins.find((item) => item.id === 'com.gcoll.modbus-tcp'))

onMounted(() => {
  consoleStore.loadPlugins()
})
</script>

<template>
  <PageHeader title="Modbus TCP 采集" description="插件详情、能力边界和使用文档。">
    <template #actions>
      <NButton secondary @click="router.push('/plugins')">
        <template #icon><ArrowLeft :size="16" /></template>
        返回
      </NButton>
      <NButton secondary :loading="consoleStore.loading" @click="consoleStore.loadPlugins">
        <template #icon><RefreshCw :size="16" /></template>
        刷新
      </NButton>
      <NButton type="primary" @click="router.push('/devices')">
        <template #icon><ListChecks :size="16" /></template>
        添加设备
      </NButton>
    </template>
  </PageHeader>

  <div class="plugin-detail-layout">
    <SectionPanel title="插件信息" subtitle="运行时和权限由插件中心管理。">
      <NDescriptions v-if="plugin" :column="2" size="small" label-placement="left" bordered>
        <NDescriptionsItem label="插件 ID">{{ plugin.id }}</NDescriptionsItem>
        <NDescriptionsItem label="状态">
          <StatusBadge :label="plugin.status === 'running' ? '运行中' : '已停止'" :kind="plugin.status === 'running' ? 'success' : 'neutral'" />
        </NDescriptionsItem>
        <NDescriptionsItem label="类型">
          <NTag size="small" type="info" :bordered="false">南向插件</NTag>
        </NDescriptionsItem>
        <NDescriptionsItem label="版本">{{ plugin.version }}</NDescriptionsItem>
        <NDescriptionsItem label="运行时">{{ plugin.runtime }}</NDescriptionsItem>
        <NDescriptionsItem label="协议">{{ plugin.protocol }}</NDescriptionsItem>
        <NDescriptionsItem label="更新时间">{{ plugin.updatedAt }}</NDescriptionsItem>
        <NDescriptionsItem label="权限">
          <span class="tag-list">
            <NTag v-for="permission in plugin.permissions" :key="permission" size="small" :bordered="false">
              {{ permission }}
            </NTag>
          </span>
        </NDescriptionsItem>
      </NDescriptions>
    </SectionPanel>

    <SectionPanel title="插件用途" class="plugin-detail-main">
      <div class="doc-block">
        <BookOpen :size="18" />
        <div>
          <strong>通过 Modbus TCP 读取和写入工业设备点位</strong>
          <p>支持线圈、离散输入、保持寄存器和输入寄存器读取；写入只允许线圈和保持寄存器。</p>
        </div>
      </div>
      <div class="doc-grid">
        <section>
          <h3>配置入口</h3>
          <p>连接参数、上报模式、调试开关和协议优化参数属于设备实例，添加或编辑设备时选择该协议后配置。</p>
        </section>
        <section>
          <h3>点位模型</h3>
          <p>通用点位表保存点位名称、地址和值类型；Modbus 扩展信息放在点位 metadata 中。</p>
        </section>
        <section>
          <h3>调试方式</h3>
          <p>设备配置页提供测试按钮，测试结果展示在当前设备的测试面板中，调试日志只保留有限窗口。</p>
        </section>
        <section>
          <h3>采集优化</h3>
          <p>正式启动采集任务时根据点位表生成读取计划，后台定期根据延迟、超时和错误调整批量长度。</p>
        </section>
      </div>
    </SectionPanel>
  </div>
</template>

<style scoped>
.plugin-detail-layout {
  display: grid;
  grid-template-columns: minmax(360px, 0.72fr) minmax(0, 1.28fr);
  gap: 8px;
}

.plugin-detail-main {
  min-height: 320px;
}

.doc-block {
  display: grid;
  grid-template-columns: auto minmax(0, 1fr);
  gap: 10px;
  align-items: start;
  padding: 12px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
  background: var(--color-bg-subtle);
}

.doc-block strong {
  display: block;
  color: var(--color-text-primary);
  font-size: 14px;
  line-height: 20px;
}

.doc-block p,
.doc-grid p {
  margin: 4px 0 0;
  color: var(--color-text-secondary);
  font-size: 13px;
  line-height: 20px;
}

.doc-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
  margin-top: 10px;
}

.doc-grid section {
  min-width: 0;
  padding: 10px;
  border: 1px solid var(--color-border);
  border-radius: 6px;
}

.doc-grid h3 {
  margin: 0;
  color: var(--color-text-primary);
  font-size: 13px;
  line-height: 19px;
}

@media (max-width: 1180px) {
  .plugin-detail-layout,
  .doc-grid {
    grid-template-columns: minmax(0, 1fr);
  }
}
</style>
