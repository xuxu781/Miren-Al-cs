<template>
  <div class="cdkeys-container">
    <div class="page-header">
      <h2>卡密管理</h2>
      <el-button type="primary" @click="showGenerateDialog">
        <el-icon><Plus /></el-icon> 生成卡密
      </el-button>
    </div>

    <el-card class="box-card" shadow="never">
      <div class="toolbar">
        <div class="search-bar">
          <el-input
            v-model="searchQuery"
            placeholder="搜索卡密"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
            class="search-input"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          
          <el-select v-model="statusFilter" placeholder="状态" clearable @change="handleSearch" class="status-select">
            <el-option label="未使用" value="0" />
            <el-option label="已使用" value="1" />
            <el-option label="已作废" value="2" />
            <el-option label="已过期" value="3" />
          </el-select>

          <el-input
            v-model="pointsFilter"
            placeholder="搜索积分额度"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
            class="points-input"
            type="number"
          />

          <el-input
            v-model="userFilter"
            placeholder="搜索使用人邮箱"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
            class="user-input"
          />
          
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button type="success" @click="handleExport" :loading="exporting">导出筛选</el-button>
          <el-button type="danger" :disabled="selectedIds.length === 0" @click="handleBatchDelete" :loading="batchDeleting">批量删除</el-button>
        </div>
      </div>

      <el-table :data="cdkeys" style="width: 100%" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="cdkey" label="卡密" min-width="250">
          <template #default="{ row }">
            <span class="font-mono">{{ row.cdkey }}</span>
            <el-button link type="primary" size="small" @click="copyCdkey(row.cdkey)" class="ml-2">复制</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="points" label="包含积分" width="120">
          <template #default="{ row }">
            <el-tag type="warning" effect="plain">{{ row.points }} 积分</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? (isExpired(row.expires_at) ? 'danger' : 'success') : (row.status === 1 ? 'info' : 'danger')">
              {{ row.status === 0 ? (isExpired(row.expires_at) ? '已过期' : '未使用') : (row.status === 1 ? '已使用' : '已作废') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="max_uses" label="使用限制" width="120">
          <template #default="{ row }">
            <span v-if="row.max_uses === 1">单次</span>
            <span v-else-if="row.max_uses === 0">无限制 (已用 {{ row.current_uses }})</span>
            <span v-else>{{ row.current_uses }} / {{ row.max_uses }} 次</span>
          </template>
        </el-table-column>
        <el-table-column prop="used_by_user_email" label="使用记录" width="160">
          <template #default="{ row }">
            <span v-if="row.max_uses === 1 && row.used_by_user_email">{{ row.used_by_user_email }}</span>
            <el-button 
              v-else-if="row.max_uses !== 1 && row.current_uses > 0" 
              link 
              type="primary" 
              size="small" 
              @click="showUsages(row.id)"
            >
              查看记录 ({{ row.current_uses }})
            </el-button>
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="expires_at" label="过期时间" width="180">
          <template #default="{ row }">
            <span v-if="row.expires_at">{{ formatDate(row.expires_at) }}</span>
            <span v-else class="text-gray-400">永久有效</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column prop="used_at" label="使用时间" width="180">
          <template #default="{ row }">
            <span v-if="row.used_at">{{ formatDate(row.used_at) }}</span>
            <span v-else class="text-gray-400">-</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm
              v-if="row.status !== 2"
              title="确定作废此卡密吗？"
              @confirm="handleVoid(row.id)"
            >
              <template #reference>
                <el-button link type="warning" size="small">作废</el-button>
              </template>
            </el-popconfirm>
            <el-popconfirm
              v-else
              title="确定启用此卡密吗？"
              @confirm="handleEnable(row.id)"
            >
              <template #reference>
                <el-button link type="success" size="small">启用</el-button>
              </template>
            </el-popconfirm>
            <el-popconfirm
              title="确定删除此卡密吗？"
              @confirm="handleDelete(row.id)"
            >
              <template #reference>
                <el-button link type="danger" size="small">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 生成卡密对话框 -->
    <el-dialog v-model="dialogVisible" title="生成卡密" width="450px">
      <el-form :model="generateForm" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="生成数量" prop="count">
          <el-input-number v-model="generateForm.count" :min="1" :max="1000" />
        </el-form-item>
        <el-form-item label="单张积分" prop="points">
          <el-select 
            v-model="generateForm.points" 
            placeholder="请选择或输入积分面额" 
            style="width: 100%"
            filterable
            allow-create
            default-first-option
          >
            <el-option label="10 积分" :value="10" />
            <el-option label="50 积分" :value="50" />
            <el-option label="100 积分" :value="100" />
            <el-option label="200 积分" :value="200" />
            <el-option label="300 积分" :value="300" />
            <el-option label="500 积分" :value="500" />
            <el-option label="1000 积分" :value="1000" />
          </el-select>
        </el-form-item>
        <el-form-item label="卡密前缀" prop="prefix">
          <el-input v-model="generateForm.prefix" placeholder="默认为 Mirenai" />
          <div class="text-xs text-gray-400 ml-2" style="line-height: 1.2; margin-top: 8px;">
            卡密生成时的前缀（系统会自动在后面添加 "-"），不填则默认使用 Mirenai
          </div>
        </el-form-item>
        <el-form-item label="多用户使用">
          <el-switch v-model="generateForm.isMultiUser" />
          <div class="text-xs text-gray-400 ml-2" style="line-height: 1.2; margin-top: 8px;">
            开启后，该卡密可被多个不同用户兑换。
          </div>
        </el-form-item>
        <el-form-item label="使用次数" v-if="generateForm.isMultiUser" prop="max_uses">
          <el-input-number v-model="generateForm.max_uses" :min="0" :max="10000" />
          <div class="text-xs text-gray-400 ml-2" style="line-height: 1.2; margin-top: 8px;">
            设置为 0 表示无限制次数。
          </div>
        </el-form-item>
        <el-form-item label="有效期至" prop="expires_at">
          <el-date-picker
            v-model="generateForm.expires_at"
            type="datetime"
            placeholder="选择过期时间"
            style="width: 100%"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
          <div class="text-xs text-gray-400 ml-2" style="line-height: 1.2; margin-top: 8px;">
            不选择表示永久有效。
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitGenerate" :loading="generating">生成</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑卡密对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑卡密" width="450px">
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="100px">
        <el-form-item label="单张积分" prop="points">
          <el-input-number v-model="editForm.points" :min="1" :max="100000" />
        </el-form-item>
        <el-form-item label="多用户使用">
          <el-switch v-model="editForm.isMultiUser" />
        </el-form-item>
        <el-form-item label="使用次数" v-if="editForm.isMultiUser" prop="max_uses">
          <el-input-number v-model="editForm.max_uses" :min="0" :max="10000" />
          <div class="text-xs text-gray-400 ml-2" style="line-height: 1.2; margin-top: 8px;">
            设置为 0 表示无限制次数。
          </div>
        </el-form-item>
        <el-form-item label="有效期至" prop="expires_at">
          <el-date-picker
            v-model="editForm.expires_at"
            type="datetime"
            placeholder="选择过期时间"
            style="width: 100%"
            value-format="YYYY-MM-DDTHH:mm:ssZ"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitEdit" :loading="editing">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 使用记录对话框 -->
    <el-dialog v-model="usagesDialogVisible" title="卡密使用记录" width="500px">
      <el-table :data="usageList" style="width: 100%" v-loading="loadingUsages">
        <el-table-column prop="email" label="使用者邮箱" min-width="200" />
        <el-table-column prop="created_at" label="使用时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="usagesDialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 生成结果展示对话框 -->
    <el-dialog v-model="resultDialogVisible" title="卡密生成成功" width="500px">
      <div class="result-container">
        <el-input
          type="textarea"
          v-model="generatedKeysText"
          :rows="10"
          readonly
          class="mb-4"
        />
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="resultDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="copyAllKeys">一键复制全部</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import axios from 'axios'
import type { FormInstance, FormRules } from 'element-plus'

const loading = ref(false)
const cdkeys = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const statusFilter = ref('')
const pointsFilter = ref('')
const userFilter = ref('')

const selectedIds = ref<number[]>([])
const exporting = ref(false)
const batchDeleting = ref(false)

const handleSelectionChange = (selection: any[]) => {
  selectedIds.value = selection.map(item => item.id)
}

const dialogVisible = ref(false)
const generating = ref(false)
const formRef = ref<FormInstance>()

const editDialogVisible = ref(false)
const editing = ref(false)
const editFormRef = ref<FormInstance>()

const usagesDialogVisible = ref(false)
const loadingUsages = ref(false)
const usageList = ref([])

const resultDialogVisible = ref(false)
const generatedKeysText = ref('')

const generateForm = reactive({
  count: 10,
  points: 100,
  prefix: 'Mirenai',
  isMultiUser: false,
  max_uses: 100,
  expires_at: ''
})

const editForm = reactive({
  id: 0,
  points: 100,
  isMultiUser: false,
  max_uses: 100,
  expires_at: ''
})

const rules = reactive<FormRules>({
  count: [{ required: true, message: '请输入生成数量', trigger: 'blur' }],
  points: [{ required: true, message: '请输入单张积分', trigger: 'blur' }]
})

const editRules = reactive<FormRules>({
  points: [{ required: true, message: '请输入单张积分', trigger: 'blur' }]
})

const fetchCdkeys = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('token')
    let url = `${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/cdkeys/list?page=${currentPage.value}&page_size=${pageSize.value}`
    if (searchQuery.value) {
      url += `&search=${encodeURIComponent(searchQuery.value)}`
    }
    if (statusFilter.value) {
      url += `&status=${statusFilter.value}`
    }
    if (pointsFilter.value) {
      url += `&points=${pointsFilter.value}`
    }
    if (userFilter.value) {
      url += `&user_email=${encodeURIComponent(userFilter.value)}`
    }
    const res = await axios.get(url, {
      headers: { Authorization: `Bearer ${token}` }
    })
    cdkeys.value = res.data.list
    total.value = res.data.total
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '获取卡密列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchCdkeys()
}

const handleExport = async () => {
  // 检查是否至少有一个筛选条件
  if (!searchQuery.value && !statusFilter.value && !pointsFilter.value && !userFilter.value) {
    ElMessage.warning('请先输入或选择筛选条件，不能导出全部数据')
    return
  }

  exporting.value = true
  try {
    const token = localStorage.getItem('token')
    let url = `${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/cdkeys/list?export=1`
    if (searchQuery.value) url += `&search=${encodeURIComponent(searchQuery.value)}`
    if (statusFilter.value) url += `&status=${statusFilter.value}`
    if (pointsFilter.value) url += `&points=${pointsFilter.value}`
    if (userFilter.value) url += `&user_email=${encodeURIComponent(userFilter.value)}`

    const res = await axios.get(url, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    const list = res.data.list || []
    if (list.length === 0) {
      ElMessage.warning('没有可导出的数据')
      return
    }

    const headers = ['ID', '卡密', '包含积分', '状态', '使用限制', '已使用次数', '最后使用邮箱', '过期时间', '创建时间']
    const csvContent = [
      headers.join(','),
      ...list.map((item: any) => {
        const statusText = item.status === 0 ? (isExpired(item.expires_at) ? '已过期' : '未使用') : (item.status === 1 ? '已使用' : '已作废')
        const limitText = item.max_uses === 1 ? '单次' : (item.max_uses === 0 ? '无限制' : item.max_uses + '次')
        const emailText = item.used_by_user_email || '-'
        const expireText = item.expires_at ? formatDate(item.expires_at) : '永久有效'
        const createText = formatDate(item.created_at)
        
        return [
          item.id,
          item.cdkey,
          item.points,
          statusText,
          limitText,
          item.current_uses,
          emailText,
          expireText,
          createText
        ].map(val => `"${val}"`).join(',')
      })
    ].join('\n')

    const blob = new Blob(['\uFEFF' + csvContent], { type: 'text/csv;charset=utf-8;' })
    const link = document.createElement('a')
    link.href = URL.createObjectURL(blob)
    link.download = `cdkeys_export_${new Date().getTime()}.csv`
    link.click()
    URL.revokeObjectURL(link.href)
    
    ElMessage.success('导出成功')
  } catch (error: any) {
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

const handleBatchDelete = () => {
  ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.value.length} 个卡密吗？`, '批量删除', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    batchDeleting.value = true
    try {
      const token = localStorage.getItem('token')
      await axios.post(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/cdkeys/batch-delete`, { ids: selectedIds.value }, {
        headers: { Authorization: `Bearer ${token}` }
      })
      ElMessage.success('批量删除成功')
      fetchCdkeys()
      // 清空选中
      selectedIds.value = []
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || '批量删除失败')
    } finally {
      batchDeleting.value = false
    }
  }).catch(() => {})
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchCdkeys()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchCdkeys()
}

const showGenerateDialog = () => {
  generateForm.count = 10
  generateForm.points = 100
  generateForm.prefix = 'Mirenai'
  generateForm.isMultiUser = false
  generateForm.max_uses = 100
  generateForm.expires_at = ''
  dialogVisible.value = true
}

const submitGenerate = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      // 确保如果是自定义输入的字符串，转换为数字
      const parsedPoints = Number(generateForm.points)
      if (isNaN(parsedPoints) || parsedPoints <= 0) {
        ElMessage.error('请输入有效的积分面额')
        return
      }

      generating.value = true
      try {
        const token = localStorage.getItem('token')
        const payload: any = {
          count: generateForm.count,
          points: parsedPoints,
          prefix: generateForm.prefix,
          max_uses: generateForm.isMultiUser ? generateForm.max_uses : 1
        }
        if (generateForm.expires_at) {
          payload.expires_at = generateForm.expires_at
        }
        const res = await axios.post(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/cdkeys/generate`, payload, {
          headers: { Authorization: `Bearer ${token}` }
        })
        ElMessage.success(res.data.message)
        dialogVisible.value = false
        
        // 显示生成的卡密结果
        if (res.data.keys && res.data.keys.length > 0) {
          generatedKeysText.value = res.data.keys.join('\n')
          resultDialogVisible.value = true
        }
        
        currentPage.value = 1
        fetchCdkeys()
      } catch (error: any) {
        ElMessage.error(error.response?.data?.error || '生成失败')
      } finally {
        generating.value = false
      }
    }
  })
}

const handleEdit = (row: any) => {
  editForm.id = row.id
  editForm.points = row.points
  editForm.isMultiUser = row.max_uses !== 1
  editForm.max_uses = row.max_uses === 1 ? 100 : row.max_uses
  editForm.expires_at = row.expires_at ? new Date(row.expires_at).toISOString() : ''
  editDialogVisible.value = true
}

const submitEdit = async () => {
  if (!editFormRef.value) return
  await editFormRef.value.validate(async (valid) => {
    if (valid) {
      editing.value = true
      try {
        const token = localStorage.getItem('token')
        const payload: any = {
          id: editForm.id,
          points: editForm.points,
          max_uses: editForm.isMultiUser ? editForm.max_uses : 1
        }
        if (editForm.expires_at) {
          payload.expires_at = editForm.expires_at
        }
        const res = await axios.post(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/cdkeys/update`, payload, {
          headers: { Authorization: `Bearer ${token}` }
        })
        ElMessage.success(res.data.message)
        editDialogVisible.value = false
        fetchCdkeys()
      } catch (error: any) {
        ElMessage.error(error.response?.data?.error || '保存失败')
      } finally {
        editing.value = false
      }
    }
  })
}

const showUsages = async (id: number) => {
  usagesDialogVisible.value = true
  loadingUsages.value = true
  try {
    const token = localStorage.getItem('token')
    const res = await axios.get(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/cdkeys/usages?cdkey_id=${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    usageList.value = res.data.list
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '获取使用记录失败')
  } finally {
    loadingUsages.value = false
  }
}

const handleVoid = async (id: number) => {
  try {
    const token = localStorage.getItem('token')
    await axios.post(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/cdkeys/void`, { id }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    ElMessage.success('作废成功')
    fetchCdkeys()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '作废失败')
  }
}

const handleDelete = async (id: number) => {
  try {
    const token = localStorage.getItem('token')
    await axios.post(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/cdkeys/delete`, { id }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    ElMessage.success('删除成功')
    fetchCdkeys()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '删除失败')
  }
}

const copyAllKeys = () => {
  navigator.clipboard.writeText(generatedKeysText.value).then(() => {
    ElMessage.success('所有卡密已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

const copyCdkey = (text: string) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制卡密到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

const handleEnable = async (id: number) => {
  try {
    const token = localStorage.getItem('token')
    await axios.post(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/cdkeys/enable`, { id }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    ElMessage.success('启用成功')
    fetchCdkeys()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '启用失败')
  }
}

const formatDate = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const isExpired = (expiresAt: string) => {
  if (!expiresAt) return false
  return new Date(expiresAt).getTime() < new Date().getTime()
}

onMounted(() => {
  fetchCdkeys()
})
</script>

<style scoped>
.cdkeys-container {
  padding: 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
  color: #161823;
}

.toolbar {
  margin-bottom: 16px;
}

.search-bar {
  display: flex;
  gap: 12px;
}

.search-input {
  width: 300px;
}

.status-select {
  width: 120px;
}

.points-input {
  width: 140px;
}

.user-input {
  width: 200px;
}

.font-mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
  background-color: #f3f4f6;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 13px;
  color: #374151;
}

.ml-2 {
  margin-left: 8px;
}

.text-gray-400 {
  color: #9ca3af;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
