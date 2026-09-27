<template>
  <div class="app-container">
    <!-- 搜索表单 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="queryParams" ref="queryFormRef" class="search-form">
        <el-form-item label="关键字" prop="keyword">
          <el-input
            v-model="queryParams.keyword"
            placeholder="请输入用户名或提示词"
            clearable
            style="width: 240px"
            @keyup.enter="handleQuery"
          >
            <template #prefix><el-icon><Search /></el-icon></template>
          </el-input>
        </el-form-item>
        <el-form-item label="任务状态" prop="status">
          <el-select
            v-model="queryParams.status"
            placeholder="请选择状态"
            clearable
            style="width: 240px"
          >
            <el-option label="全部" value="" />
            <el-option label="处理中" :value="0" />
            <el-option label="成功" :value="1" />
            <el-option label="部分成功" :value="3" />
            <el-option label="失败" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleQuery">搜索</el-button>
          <el-button :icon="Refresh" @click="resetQuery">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 表格区域 -->
    <el-card shadow="never" class="table-card">
      <div class="table-header">
        <span class="table-title">任务列表</span>
        <div class="table-actions">
          <el-tooltip content="刷新" placement="top">
            <el-button circle :icon="Refresh" @click="fetchTasks" />
          </el-tooltip>
        </div>
      </div>

      <el-table
        v-loading="loading"
        :data="tasks"
        border
        stripe
        style="width: 100%"
        :header-cell-style="{ background: '#f8f8f9', color: '#515a6e', fontWeight: 'bold' }"
      >
        <el-table-column label="任务ID" prop="id" width="80" align="center" />
        <el-table-column label="用户名" prop="username" width="120" align="center" />
        <el-table-column label="调用模型" prop="model" width="160" align="center">
          <template #default="scope">
            <div v-if="scope.row.model || scope.row.series_id" style="display: flex; flex-direction: column; align-items: center; gap: 4px;">
              <el-tag size="small" type="primary" effect="light" v-if="scope.row.series_id">
                {{ scope.row.series_id }}
              </el-tag>
              <el-tag size="small" type="danger" effect="light" v-if="scope.row.model && scope.row.model !== scope.row.series_id">
                {{ scope.row.model }}
              </el-tag>
              <el-tag size="small" type="warning" effect="light" v-if="scope.row.channel_name">
                渠道: {{ scope.row.channel_name }}
              </el-tag>
            </div>
            <span v-else class="empty-text">--</span>
          </template>
        </el-table-column>
        <el-table-column label="提示词" prop="prompt" min-width="250" show-overflow-tooltip />
        <el-table-column label="耗时" width="90" align="center">
          <template #default="scope">
            <span v-if="scope.row.status !== 0 && scope.row.created_at && scope.row.updated_at" style="color: #606266; font-family: monospace;">
              {{ getDuration(scope.row.created_at, scope.row.updated_at) }}
            </span>
            <span v-else-if="scope.row.status === 0" style="color: #909399; font-style: italic;">
              计算中...
            </span>
            <span v-else class="empty-text">--</span>
          </template>
        </el-table-column>
        <el-table-column label="生成尺寸" prop="size" width="100" align="center">
          <template #default="scope">
            <el-tag size="small" type="info">{{ scope.row.size || '默认' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="分辨率" prop="resolution" width="100" align="center">
          <template #default="scope">
            <el-tag size="small" type="warning" v-if="scope.row.resolution || getResolution(scope.row.log_content) !== '--'">
              {{ scope.row.resolution || getResolution(scope.row.log_content) }}
            </el-tag>
            <span v-else class="empty-text">--</span>
          </template>
        </el-table-column>
        <el-table-column label="生成明细" width="140" align="center">
          <template #default="scope">
            <div style="display: flex; flex-direction: column; gap: 4px; font-size: 12px; align-items: center;">
              <span>数量: {{ scope.row.num_images || 1 }} 张</span>
              <span style="color: #f56c6c;">消耗: {{ calculateConsumedPoints(scope.row) }} 积分</span>
              <span v-if="scope.row.status === 2 || scope.row.status === 3" style="color: #67c23a;">退还: {{ calculateRefundedPoints(scope.row) }} 积分</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="结果图片" width="160" align="center">
          <template #default="scope">
            <div v-if="(scope.row.image_url && scope.row.image_url !== '[]') || scope.row.status === 0 || scope.row.status === 2 || scope.row.status === 3" class="table-images-container">
              <el-image
                v-for="(img, idx) in parseImageUrls(scope.row.image_url)"
                :key="'img-'+idx"
                :src="formatImageUrl(img)"
                :preview-src-list="parseImageUrls(scope.row.image_url).map(u => formatImageUrl(u))"
                :initial-index="idx"
                fit="cover"
                class="table-image"
                preview-teleported
              >
                <template #error>
                  <div class="image-error">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
              <!-- 失败占位符 -->
              <template v-if="scope.row.status === 2 || scope.row.status === 3">
                <div
                  v-for="n in Math.max(0, (scope.row.num_images || 1) - parseImageUrls(scope.row.image_url).length)"
                  :key="'failed-'+n"
                  class="table-image"
                  style="background-color: #fef2f2; border: 1px dashed #fca5a5; display: flex; flex-direction: column; justify-content: center; align-items: center; border-radius: 4px; box-sizing: border-box; padding: 2px;"
                >
                  <el-icon color="#f87171" :size="14"><CircleClose /></el-icon>
                  <span style="font-size: 10px; color: #ef4444; line-height: 1.1; margin-top: 2px; text-align: center;">失败</span>
                </div>
              </template>
              <!-- 生成中占位符 -->
              <template v-if="scope.row.status === 0">
                <div
                  v-for="n in Math.max(0, (scope.row.num_images || 1) - parseImageUrls(scope.row.image_url).length)"
                  :key="'loading-'+n"
                  class="table-image"
                  style="background-color: #f5f7fa; border: 1px dashed #dcdfe6; display: flex; flex-direction: column; justify-content: center; align-items: center; border-radius: 4px; box-sizing: border-box; padding: 2px;"
                >
                  <el-icon class="is-loading" color="#909399" :size="14"><Loading /></el-icon>
                  <span style="font-size: 10px; color: #909399; line-height: 1.1; margin-top: 2px; text-align: center;">生成中</span>
                </div>
              </template>
            </div>
            <span v-else class="empty-text">--</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" prop="status" width="120" align="center">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              <div style="display: flex; align-items: center; gap: 4px;">
                <span class="status-dot" :class="getStatusClass(scope.row.status)"></span>
                {{ getStatusText(scope.row.status) }}
              </div>
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" prop="created_at" width="170" align="center">
          <template #default="scope">
            {{ formatDate(scope.row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" align="center" fixed="right">
          <template #default="scope">
            <el-button
              link
              type="primary"
              :icon="Document"
              @click="viewDetails(scope.row)"
            >
              任务详情
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页区域 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="queryParams.page"
          v-model:page-size="queryParams.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>

    <!-- 任务详情抽屉 -->
    <el-drawer
      v-model="detailsDrawerVisible"
      title="任务详细信息"
      size="800px"
      destroy-on-close
    >
      <div class="task-details-container" v-if="currentTask" v-loading="logLoading">
        <el-tabs v-model="activeDetailTab" class="custom-tabs">
          <el-tab-pane label="基本信息" name="info">
            <el-descriptions :column="2" border class="details-desc">
              <el-descriptions-item label="任务 ID" width="100">{{ currentTask.id }}</el-descriptions-item>
              <el-descriptions-item label="用户名">{{ currentTask.username }}</el-descriptions-item>
              
              <el-descriptions-item label="任务状态">
                <el-tag :type="getStatusType(currentTask.status)">
                  <div style="display: flex; align-items: center; gap: 4px;">
                    <span class="status-dot" :class="getStatusClass(currentTask.status)"></span>
                    {{ getStatusText(currentTask.status) }}
                  </div>
                </el-tag>
              </el-descriptions-item>
              
              <el-descriptions-item label="调用模型">
                <div v-if="currentTask.model || currentTask.series_id" style="display: flex; flex-wrap: wrap; gap: 8px;">
                  <el-tag size="small" type="primary" effect="light" v-if="currentTask.series_id">
                    {{ currentTask.series_id }}
                  </el-tag>
                  <el-tag size="small" type="danger" effect="light" v-if="currentTask.model && currentTask.model !== currentTask.series_id">
                    {{ currentTask.model }}
                  </el-tag>
                </div>
                <span v-else class="empty-text">--</span>
              </el-descriptions-item>

              <el-descriptions-item label="生成尺寸">
                <el-tag size="small" type="info">{{ currentTask.size || '默认' }}</el-tag>
              </el-descriptions-item>

              <el-descriptions-item label="分辨率">
                <el-tag size="small" type="warning" v-if="currentTask.resolution || getResolution(rawLogContent || currentTask.log_content) !== '--'">
                  {{ currentTask.resolution || getResolution(rawLogContent || currentTask.log_content) }}
                </el-tag>
                <span v-else class="empty-text">--</span>
              </el-descriptions-item>

              <el-descriptions-item label="创建时间">{{ formatDate(currentTask.created_at) }}</el-descriptions-item>
              <el-descriptions-item label="更新时间">{{ formatDate(currentTask.updated_at) }}</el-descriptions-item>
              
              <el-descriptions-item label="任务耗时">
                <span v-if="currentTask.status !== 0 && currentTask.created_at && currentTask.updated_at" style="color: #606266; font-weight: bold; font-family: monospace;">
                  {{ getDuration(currentTask.created_at, currentTask.updated_at) }}
                </span>
                <span v-else-if="currentTask.status === 0" style="color: #909399; font-style: italic;">计算中...</span>
                <span v-else class="empty-text">--</span>
              </el-descriptions-item>

              <el-descriptions-item label="生成明细" :span="2">
                <div style="display: flex; gap: 24px; font-weight: bold;">
                  <span>请求张数: <span style="color: #409eff;">{{ currentTask.num_images || 1 }} 张</span></span>
                  <span>消耗积分: <span style="color: #f56c6c;">{{ calculateConsumedPoints(currentTask) }} 积分</span></span>
                  <span v-if="currentTask.status === 2 || currentTask.status === 3">退还积分: <span style="color: #67c23a;">{{ calculateRefundedPoints(currentTask) }} 积分</span></span>
                </div>
              </el-descriptions-item>

              <el-descriptions-item label="会话 ID" :span="2">
                {{ currentTask.session_id || '--' }}
              </el-descriptions-item>

              <el-descriptions-item label="提示词 (Prompt)" :span="2">
                <div class="prompt-box">{{ currentTask.prompt || '--' }}</div>
              </el-descriptions-item>

              <el-descriptions-item label="参考图片" :span="2">
                <div class="images-grid" v-if="currentTask.reference_image && currentTask.reference_image !== '[]'">
                  <div class="image-wrapper" v-for="(img, idx) in parseImageUrls(currentTask.reference_image)" :key="idx">
                    <el-image
                      :src="formatImageUrl(img)"
                      :preview-src-list="parseImageUrls(currentTask.reference_image).map(u => formatImageUrl(u))"
                      :initial-index="idx"
                      fit="contain"
                      class="detail-image"
                      preview-teleported
                    >
                      <template #error>
                        <div class="image-error"><el-icon><Picture /></el-icon></div>
                      </template>
                    </el-image>
                  </div>
                </div>
                <span v-else class="empty-text">无参考图</span>
              </el-descriptions-item>

              <el-descriptions-item label="结果图片" :span="2">
                <div class="images-grid" v-if="(currentTask.image_url && currentTask.image_url !== '[]') || currentTask.status === 0 || currentTask.status === 2 || currentTask.status === 3">
                  <div class="image-wrapper" v-for="(img, idx) in parseImageUrls(currentTask.image_url)" :key="'res-'+idx">
                    <el-image
                      :src="formatImageUrl(img)"
                      :preview-src-list="parseImageUrls(currentTask.image_url).map(u => formatImageUrl(u))"
                      :initial-index="idx"
                      fit="contain"
                      class="detail-image"
                      preview-teleported
                    >
                      <template #error>
                        <div class="image-error"><el-icon><Picture /></el-icon></div>
                      </template>
                    </el-image>
                  </div>
                  <!-- 失败占位符 -->
                  <template v-if="currentTask.status === 2 || currentTask.status === 3">
                    <div class="image-wrapper" v-for="n in Math.max(0, (currentTask.num_images || 1) - parseImageUrls(currentTask.image_url).length)" :key="'res-failed-'+n">
                      <div class="detail-image" style="background-color: #fef2f2; border: 1px dashed #fca5a5; display: flex; flex-direction: column; justify-content: center; align-items: center; border-radius: 4px; box-sizing: border-box; width: 100%; height: 100%; min-height: 120px; padding: 8px;">
                        <el-icon color="#f87171" :size="24"><CircleClose /></el-icon>
                        <span style="font-size: 12px; color: #ef4444; line-height: 1.2; margin-top: 8px; text-align: center;">生成失败已退款<br/><br/><span style="opacity: 0.8;">{{ getErrorMessage(currentTask.log_content, currentTask.status) }}</span></span>
                      </div>
                    </div>
                  </template>
                  <!-- 生成中占位符 -->
                  <template v-if="currentTask.status === 0">
                    <div class="image-wrapper" v-for="n in Math.max(0, (currentTask.num_images || 1) - parseImageUrls(currentTask.image_url).length)" :key="'res-loading-'+n">
                      <div class="detail-image" style="background-color: #f5f7fa; border: 1px dashed #dcdfe6; display: flex; flex-direction: column; justify-content: center; align-items: center; border-radius: 4px; box-sizing: border-box; width: 100%; height: 100%; min-height: 120px;">
                        <el-icon class="is-loading" color="#909399" :size="24"><Loading /></el-icon>
                        <span style="font-size: 12px; color: #909399; line-height: 1.2; margin-top: 8px; text-align: center;">生成中</span>
                      </div>
                    </div>
                  </template>
                </div>
                <span v-else class="empty-text">无结果图</span>
              </el-descriptions-item>

              <el-descriptions-item label="错误信息" :span="2" v-if="currentTask.status === 2 || currentTask.status === 3">
                <div class="error-box" style="color: #f56c6c; font-size: 13px; background-color: #fef0f0; padding: 12px; border-radius: 4px; border: 1px solid #fde2e2; white-space: pre-wrap; word-break: break-all;">
                  <div style="display: flex; align-items: flex-start; gap: 8px;">
                    <el-icon style="margin-top: 2px;"><Warning /></el-icon>
                    <span>{{ getErrorMessage(currentTask.log_content, currentTask.status) }}</span>
                  </div>
                </div>
              </el-descriptions-item>
            </el-descriptions>
          </el-tab-pane>
          
          <el-tab-pane label="接口完整数据" name="json">
            <div class="json-viewer">
              <pre>{{ JSON.stringify(currentTask, null, 2) }}</pre>
            </div>
          </el-tab-pane>

          <el-tab-pane label="执行日志" name="logs">
            <div class="log-box">
              <pre v-if="rawLogContent">{{ rawLogContent }}</pre>
              <span v-else class="empty-text">暂无日志内容</span>
            </div>
          </el-tab-pane>
        </el-tabs>
      </div>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { Search, Refresh, Picture, Document, CircleClose, Loading } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

// 查询参数
const queryParams = reactive({
  page: 1,
  pageSize: 10,
  keyword: '',
  status: '' as number | ''
})

const tasks = ref<any[]>([])
const models = ref<any[]>([])
const total = ref(0)
const loading = ref(false)

const getErrorMessage = (logContent?: string, status?: number) => {
  if (!logContent) return '未知错误'
  
  try {
    const errorBodyRegex = /Upstream API returned error status.*?\nBody:\s*({.*})/i;
    const match = logContent.match(errorBodyRegex);
    if (match && match[1]) {
      const errorJson = JSON.parse(match[1]);
      if (errorJson?.error?.message) {
        const msg = errorJson.error.message.toLowerCase();
        if (msg.includes('model_not_found') || msg.includes('no available channel')) {
          return '当前模型通道不可用或正在维护，请稍后重试'
        }
        if (msg.includes('invalid token') || msg.includes('unauthorized')) {
          return '接口鉴权失败，请联系管理员检查配置'
        }
        // 检查是否有上游透传的原始文本消息
        const upstreamText = errorJson.upstream_text_message || '';

        if (msg.includes('safety') || msg.includes('violation') || upstreamText.includes('安全审核') || upstreamText.includes('色情') || upstreamText.includes('暴力')) {
           return status === 3 ? '该张图片触发安全过滤' : '提示词可能包含违规或不支持的内容，请修改后重试'
        }
        
        // 如果有详细的上游透传消息，优先展示透传消息
        if (upstreamText) {
          return `${status === 3 ? '部分生成失败' : '生成失败'}: ${upstreamText}`
        }

        return `${status === 3 ? '部分生成失败' : '生成失败'}: ${errorJson.error.message}`
      }
    }
  } catch (e) {
  }

  const lower = logContent.toLowerCase()
  const safeText = lower.replace(/"prompt"\s*:\s*".*?"/g, '').replace(/"prompt"\s*:/g, '')
  
  if (safeText.includes('safety') || safeText.includes('policy') || safeText.includes('violation') || 
      safeText.includes('blocked') || safeText.includes('nsfw') || safeText.includes('敏感') || 
      safeText.includes('违规') || safeText.includes('prompt') || safeText.includes('bad request') || 
      safeText.includes('status: 400') || safeText.includes('error status (task 1): 400')) {
    return status === 3 ? '该张图片生成失败(触发过滤或网络波动)' : '提示词可能包含违规或不支持的内容，请修改后重试'
  }
  
  if (safeText.includes('invalid token') || safeText.includes('unauthorized') || safeText.includes('balance') || 
      safeText.includes('insufficient') || safeText.includes('quota') || safeText.includes('api_key') || 
      safeText.includes('key') || safeText.includes('401') || safeText.includes('402')) {
    return status === 3 ? '接口不稳定，部分生成失败' : '接口请求失败或系统繁忙，请稍后重试'
  }
  
  if (safeText.includes('timeout') || safeText.includes('network') || safeText.includes('中断') || safeText.includes('超时')) {
    return status === 3 ? '部分图片生成超时' : '生成超时或意外中断，请检查网络或稍后重试'
  }
  
  return status === 3 ? '部分图片生成失败' : '生成失败，请稍后重试'
}

// 获取模型列表用于计算积分
const fetchModels = async () => {
  try {
    const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    const res = await axios.get(`${apiUrl}/api/admin/upstreams/list`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (res.data && res.data.success && res.data.data) {
      models.value = res.data.data
    }
  } catch (error) {
    console.error('获取模型列表失败', error)
  }
}

// 获取某分辨率的计费配置
const getTierConfigForRes = (modelInfo: any, resName: string, key: string, defaultValue: any) => {
  if (!modelInfo || !modelInfo.resolution_configs) return defaultValue
  try {
    const configs = JSON.parse(modelInfo.resolution_configs)
    if (configs && configs[resName] && configs[resName][key] !== undefined) {
      return configs[resName][key]
    }
  } catch (e) {
    // 解析失败
  }
  return defaultValue
}

// 计算任务消耗的积分
const calculateConsumedPoints = (task: any) => {
  const numImages = task.num_images || 1
  const targetId = task.series_id || task.model
  if (!targetId) return parseFloat((numImages).toFixed(2))
  
  const model = models.value.find(m => m.series_id === targetId || m.logical_model === targetId || m.name === targetId)
  if (!model) return parseFloat((numImages).toFixed(2))
  
  const creditsPerImage = getTierConfigForRes(model, task.resolution, 'credits_per_image', 1.0)
  return parseFloat((Number(creditsPerImage) * numImages).toFixed(2))
}

// 计算任务退还的积分
const calculateRefundedPoints = (task: any) => {
  if (task.status !== 2 && task.status !== 3) return '0.00'
  
  const logContent = task.log_content || ''
  
  // 匹配退还日志中的失败张数 (新版部分失败日志: Partially generated X images. Y failed.)
  const partialMatch = logContent.match(/Partially generated \d+ images\. (\d+) failed\./)
  if (partialMatch && partialMatch[1]) {
    const failedCount = parseInt(partialMatch[1], 10)
    const targetId = task.series_id || task.model
    const model = models.value.find(m => m.series_id === targetId || m.logical_model === targetId || m.name === targetId)
    if (model) {
      const creditsPerImage = getTierConfigForRes(model, task.resolution, 'credits_per_image', 1.0)
      return parseFloat((Number(creditsPerImage) * failedCount).toFixed(2))
    }
    return parseFloat((failedCount).toFixed(2)) // 如果找不到模型，默认退还失败张数对应的默认积分（1张=1分）
  }
  
  // 如果是状态 2（全部失败），退还全部消耗积分
  if (task.status === 2) {
    return calculateConsumedPoints(task)
  }
  
  return '0.00'
}

// 详情抽屉相关
const detailsDrawerVisible = ref(false)
const logLoading = ref(false)
const rawLogContent = ref('')
const currentTask = ref<any>(null)
const activeDetailTab = ref('info')

let pollingTimer: ReturnType<typeof setInterval> | null = null

const startPolling = () => {
  if (pollingTimer) return
  pollingTimer = setInterval(() => {
    const hasGenerating = tasks.value.some(t => t.status === 0)
    if (hasGenerating) {
      fetchTasksSilently()
    }
  }, 3000)
}

const stopPolling = () => {
  if (pollingTimer) {
    clearInterval(pollingTimer)
    pollingTimer = null
  }
}

const fetchTasksSilently = async () => {
  try {
    const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    const res = await axios.get(`${apiUrl}/api/admin/tasks/list`, {
      params: {
        page: queryParams.page,
        page_size: queryParams.pageSize,
        keyword: queryParams.keyword,
        status: queryParams.status !== '' ? queryParams.status : undefined
      },
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (!res.data.error) {
      tasks.value = res.data.list || []
      total.value = res.data.total || 0
      
      // Update details drawer if open and task is still generating
      if (detailsDrawerVisible.value && currentTask.value && currentTask.value.status === 0) {
        const updatedTask = tasks.value.find(t => t.id === currentTask.value.id)
        if (updatedTask) {
          currentTask.value = updatedTask
          // Silently fetch logs to update log tab
          fetchLogsSilently(updatedTask.id)
        }
      }
    }
  } catch (error) {
    // Ignore error in silent fetch
  }
}

const fetchLogsSilently = async (taskId: number) => {
  try {
    const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    const res = await axios.get(`${apiUrl}/api/admin/tasks/logs`, {
      params: { task_id: taskId },
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })
    if (!res.data.error) {
      rawLogContent.value = res.data.raw || currentTask.value?.log_content || ''
    }
  } catch (error) {
    // Ignore error
  }
}

// 获取任务列表
const fetchTasks = async () => {
  loading.value = true
  try {
    const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    const res = await axios.get(`${apiUrl}/api/admin/tasks/list`, {
      params: {
        page: queryParams.page,
        page_size: queryParams.pageSize,
        keyword: queryParams.keyword,
        status: queryParams.status !== '' ? queryParams.status : undefined
      },
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (res.data.error) {
      ElMessage.error(res.data.error)
    } else {
      tasks.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '获取任务列表失败')
  } finally {
    loading.value = false
  }
}

const handleQuery = () => {
  queryParams.page = 1
  fetchTasks()
}

const resetQuery = () => {
  queryParams.keyword = ''
  queryParams.status = ''
  handleQuery()
}

const handleSizeChange = (val: number) => {
  queryParams.pageSize = val
  fetchTasks()
}

const handleCurrentChange = (val: number) => {
  queryParams.page = val
  fetchTasks()
}

// 查看详情
const viewDetails = async (row: any) => {
  currentTask.value = row
  detailsDrawerVisible.value = true
  logLoading.value = true
  rawLogContent.value = ''
  activeDetailTab.value = 'info'
  
  try {
    const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    const res = await axios.get(`${apiUrl}/api/admin/tasks/logs`, {
      params: { task_id: row.id },
      headers: {
        Authorization: `Bearer ${localStorage.getItem('token')}`
      }
    })
    
    if (res.data.error) {
      ElMessage.error(res.data.error)
    } else {
      rawLogContent.value = res.data.raw || row.log_content || ''
    }
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '获取日志失败')
    rawLogContent.value = row.log_content || ''
  } finally {
    logLoading.value = false
  }
}

// 状态工具函数
const getStatusType = (status: number) => {
  switch (status) {
    case 0: return 'primary'
    case 1: return 'success'
    case 2: return 'danger'
    case 3: return 'warning'
    default: return 'info'
  }
}

const getStatusClass = (status: number) => {
  switch (status) {
    case 0: return 'dot-primary'
    case 1: return 'dot-success'
    case 2: return 'dot-danger'
    case 3: return 'dot-warning'
    default: return 'dot-info'
  }
}

const getStatusText = (status: number) => {
  switch (status) {
    case 0: return '处理中'
    case 1: return '成功'
    case 2: return '失败'
    case 3: return '部分成功'
    default: return '未知'
  }
}

// 时间格式化
const formatDate = (dateStr: string) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return '-'
  const pad = (n: number) => (n < 10 ? '0' + n : n)
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

// 计算耗时
const getDuration = (startStr: string, endStr: string) => {
  if (!startStr || !endStr) return '--'
  const start = new Date(startStr).getTime()
  const end = new Date(endStr).getTime()
  if (isNaN(start) || isNaN(end)) return '--'
  
  const diff = Math.max(0, end - start)
  const seconds = Math.floor(diff / 1000)
  
  if (seconds < 60) {
    return `${seconds}s`
  }
  const minutes = Math.floor(seconds / 60)
  const remainSeconds = seconds % 60
  return `${minutes}m ${remainSeconds}s`
}

// 解析分辨率
const getResolution = (logStr: string) => {
  if (!logStr) return '--'
  const match = logStr.match(/"resolution"\s*:\s*"([^"]+)"/)
  if (match && match[1]) {
    return match[1]
  }
  return '--'
}

// 图片处理
const parseImageUrls = (imgData: string) => {
  if (!imgData) return []
  try {
    const arr = JSON.parse(imgData)
    if (Array.isArray(arr)) return arr
  } catch (e) {}
  return [imgData]
}

const formatImageUrl = (img: string) => {
  if (!img) return ''
  if (!img.startsWith('[') && !img.startsWith('{')) {
    if (img.startsWith('http://localhost:8080/uploads/')) {
      img = img.replace('http://localhost:8080', '')
    }
    if (img.length > 1000 && !img.startsWith('http') && !img.startsWith('data:')) {
      return 'data:image/png;base64,' + img
    }
    if (img.startsWith('/uploads/')) {
      const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
      return apiUrl + img
    }
    return img
  }
  try {
    const arr = JSON.parse(img)
    if (Array.isArray(arr) && arr.length > 0) {
      let firstImg = arr[0]
      if (firstImg.startsWith('http://localhost:8080/uploads/')) {
        firstImg = firstImg.replace('http://localhost:8080', '')
      }
      if (firstImg.startsWith('/uploads/')) {
        const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
        return apiUrl + firstImg
      }
      if (firstImg.length > 1000 && !firstImg.startsWith('http') && !firstImg.startsWith('data:')) {
        return 'data:image/png;base64,' + firstImg
      }
      return firstImg
    }
  } catch (e) {}
  if (img.startsWith('http://localhost:8080/uploads/')) {
    img = img.replace('http://localhost:8080', '')
  }
  if (img.length > 1000 && !img.startsWith('http') && !img.startsWith('data:')) {
    return 'data:image/png;base64,' + img
  }
  if (img.startsWith('/uploads/')) {
    const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    return apiUrl + img
  }
  return img
}

onMounted(() => {
  fetchModels()
  fetchTasks()
  startPolling()
})

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.app-container {
  padding: 20px;
  background-color: #f0f2f5;
  min-height: calc(100vh - 84px);
}

.search-card {
  margin-bottom: 16px;
  border-radius: 4px;
}

.search-form .el-form-item {
  margin-bottom: 0;
}

.table-card {
  border-radius: 4px;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.table-title {
  font-size: 16px;
  font-weight: 600;
  color: #1f2f3d;
}

.table-actions {
  display: flex;
  gap: 12px;
}

.table-images-container {
  display: flex;
  gap: 4px;
  flex-wrap: wrap;
  justify-content: center;
}

.table-image {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  border: 1px solid #ebeef5;
}

.image-error {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  background: #f5f7fa;
  color: #909399;
  font-size: 20px;
}

/* 状态圆点样式 */
.status-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  border-radius: 50%;
}
.dot-primary { background-color: #409eff; }
.dot-success { background-color: #67c23a; }
.dot-danger { background-color: #f56c6c; }
.dot-warning { background-color: #e6a23c; }
.dot-info { background-color: #909399; }

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

/* 抽屉 & 详情样式 */
:deep(.el-drawer__header) {
  margin-bottom: 0;
  padding-bottom: 16px;
  border-bottom: 1px solid #ebeef5;
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.task-details-container {
  padding: 10px 0;
}

.details-desc {
  --el-descriptions-item-bordered-label-background: #f8f8f9;
}

.prompt-box {
  background-color: #f5f7fa;
  padding: 10px 14px;
  border-radius: 4px;
  color: #606266;
  line-height: 1.5;
  word-break: break-all;
  white-space: pre-wrap;
  font-size: 13px;
}

.images-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 12px;
  width: 100%;
}

.image-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
  aspect-ratio: 1 / 1;
  background: #f5f7fa;
  border-radius: 4px;
  padding: 8px;
  overflow: hidden;
}

.detail-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  border-radius: 4px;
}

.log-box {
  background-color: #1e1e1e;
  color: #a6e22e;
  padding: 12px 16px;
  border-radius: 6px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  max-height: 300px;
  overflow-y: auto;
}

.log-box pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}

.json-viewer {
  background-color: #f5f7fa;
  padding: 16px;
  border-radius: 6px;
  border: 1px solid #e4e7ed;
  max-height: 500px;
  overflow-y: auto;
}

.json-viewer pre {
  margin: 0;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 13px;
  color: #303133;
  white-space: pre-wrap;
  word-break: break-all;
}

.empty-text {
  color: #909399;
  font-size: 13px;
  font-style: italic;
}
</style>
