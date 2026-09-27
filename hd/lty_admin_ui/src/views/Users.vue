<template>
  <div class="admin-container fade-in">
    <div class="page-header">
      <div class="header-info">
        <div class="icon-wrapper">
          <el-icon :size="24" color="#409eff"><User /></el-icon>
        </div>
        <div class="header-title">
          <h2>用户管理</h2>
          <span class="subtitle">管理平台普通用户，修改信息及状态控制</span>
        </div>
      </div>
      <div class="header-actions">
        <el-input
          v-model="searchQuery"
          placeholder="搜索用户名/邮箱/手机号..."
          :prefix-icon="Search"
          clearable
          @clear="handleSearch"
          @keyup.enter="handleSearch"
          class="search-input"
        >
          <template #append>
            <el-button :icon="Search" @click="handleSearch" />
          </template>
        </el-input>
        <el-button type="primary" class="add-btn" @click="showAddDialog = true">
          <el-icon class="el-icon--left"><Plus /></el-icon>
          添加用户
        </el-button>
      </div>
    </div>

    <el-card shadow="never" class="admin-card">
      <el-table 
        :data="userList" 
        v-loading="loading" 
        style="width: 100%" 
        :header-cell-style="{ background: '#f8fafc', color: '#475569', fontWeight: '600', height: '54px' }"
        row-class-name="custom-table-row"
      >
        <el-table-column prop="id" label="ID" width="100" align="center">
          <template #default="scope">
            <div class="id-badge">#{{ scope.row.id }}</div>
          </template>
        </el-table-column>
        
        <el-table-column prop="username" label="用户名" min-width="160">
          <template #default="scope">
            <div class="user-info">
              <el-avatar :size="36" class="user-avatar" :style="{ background: getAvatarColor(scope.row.username) }">
                {{ scope.row.username.charAt(0).toUpperCase() }}
              </el-avatar>
              <div class="user-details">
                <span class="username">{{ scope.row.username }}</span>
              </div>
            </div>
          </template>
        </el-table-column>

        <el-table-column prop="email" label="邮箱" min-width="180">
          <template #default="scope">
            <div class="info-cell" v-if="scope.row.email">
              <el-icon><Message /></el-icon>
              <span>{{ scope.row.email }}</span>
            </div>
            <span v-else class="empty-text">-</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="phone" label="手机号" min-width="140">
          <template #default="scope">
            <div class="info-cell" v-if="scope.row.phone">
              <el-icon><Phone /></el-icon>
              <span>{{ scope.row.phone }}</span>
            </div>
            <span v-else class="empty-text">-</span>
          </template>
        </el-table-column>
        
        <el-table-column prop="points" label="积分" min-width="120">
          <template #default="scope">
            <div class="info-cell">
              <el-icon><Coin /></el-icon>
              <span style="font-weight: 600; color: #f59e0b;">{{ parseFloat(Number(scope.row.points || 0).toFixed(2)) }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="status" label="状态" width="120" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'" size="small" effect="light" class="status-tag">
              <span class="status-dot" :class="scope.row.status === 1 ? 'dot-success' : 'dot-danger'"></span>
              {{ scope.row.status === 1 ? '正常' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="created_at" label="注册时间" min-width="180" align="center">
          <template #default="scope">
            <div class="time-cell">
              <el-icon><Calendar /></el-icon>
              <span>{{ new Date(scope.row.created_at).toLocaleString() }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="260" align="center" fixed="right">
          <template #default="scope">
            <div class="action-btns">
              <el-tooltip content="积分记录" placement="top">
                <el-button type="info" link @click="openPointsRecordDialog(scope.row)" class="action-btn">
                  <el-icon :size="16"><List /></el-icon>
                </el-button>
              </el-tooltip>

              <el-tooltip content="编辑用户信息" placement="top">
                <el-button type="primary" link @click="openEditDialog(scope.row)" class="action-btn">
                  <el-icon :size="16"><EditPen /></el-icon>
                </el-button>
              </el-tooltip>
              
              <el-tooltip content="删除用户" placement="top">
                <div style="display: inline-block;">
                  <el-popconfirm 
                    title="确定要删除该用户吗？此操作不可恢复！" 
                    confirm-button-text="确定删除"
                    cancel-button-text="取消"
                    confirm-button-type="danger"
                    :icon="Warning"
                    icon-color="#f56c6c"
                    @confirm="handleDelete(scope.row)"
                  >
                    <template #reference>
                      <el-button type="danger" link class="action-btn danger-btn">
                        <el-icon :size="16"><Delete /></el-icon>
                      </el-button>
                    </template>
                  </el-popconfirm>
                </div>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>
        
        <template #empty>
          <el-empty description="暂无用户数据" :image-size="120" />
        </template>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :background="true"
          layout="total, sizes, prev, pager, next, jumper"
          :total="total"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 添加用户对话框 -->
    <el-dialog v-model="showAddDialog" title="添加新用户" width="500px" destroy-on-close class="custom-dialog">
      <div class="dialog-desc">填写基本信息以添加一个新的普通用户。</div>
      <el-form ref="addFormRef" :model="addForm" :rules="rules" label-width="80px" size="large" label-position="left">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="addForm.username" placeholder="请输入用户名（留空则默认使用邮箱）" />
        </el-form-item>
        <el-form-item label="初始密码" prop="password">
          <el-input v-model="addForm.password" type="password" placeholder="请输入初始密码 (至少6位)" show-password />
        </el-form-item>
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="addForm.email" placeholder="请输入邮箱地址" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="addForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="积分" prop="points">
          <el-input-number v-model="addForm.points" :min="0" :step="10" :precision="2" placeholder="请输入初始积分" />
        </el-form-item>
        <el-form-item label="账号状态" prop="status">
          <el-switch v-model="addForm.status" :active-value="1" :inactive-value="0" active-text="正常" inactive-text="禁用" inline-prompt />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddDialog = false">取 消</el-button>
          <el-button type="primary" @click="handleAdd" :loading="submitLoading">确认添加</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑用户对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑用户信息" width="500px" destroy-on-close class="custom-dialog">
      <div class="dialog-desc">修改用户信息，如果不修改密码请留空。</div>
      <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="80px" size="large" label-position="left">
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="editForm.email" disabled placeholder="请输入邮箱地址" />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="editForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="新密码" prop="password">
          <el-input v-model="editForm.password" type="password" placeholder="留空则不修改密码" show-password />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="editForm.phone" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="当前积分">
          <span style="font-weight: bold; color: #f59e0b;">{{ originalPoints }}</span>
        </el-form-item>
        <el-form-item label="操作">
          <el-radio-group v-model="editForm.pointAction" @change="handlePointActionChange">
            <el-radio label="none">不修改</el-radio>
            <el-radio label="add">增加积分</el-radio>
            <el-radio label="subtract" @click.native="handleSubtractClick">扣除积分</el-radio>
          </el-radio-group>
          <div style="font-size: 12px; color: #f56c6c; margin-left: 10px; line-height: 1;" v-if="editForm.pointAction === 'subtract' && originalPoints <= 0">
            该用户当前积分为 0，无法扣除
          </div>
        </el-form-item>
        <el-form-item label="变动数量" prop="pointsChange" v-if="editForm.pointAction !== 'none'">
          <el-input-number v-model="editForm.pointsChange" :min="0.01" :max="editForm.pointAction === 'subtract' ? originalPoints : 99999" :step="1" :precision="2" placeholder="请输入变动数量" />
          <div style="font-size: 12px; color: #f56c6c; margin-left: 10px; line-height: 1;" v-if="editForm.pointAction === 'subtract'">
            最多只能扣除 {{ originalPoints }} 积分
          </div>
        </el-form-item>
        <el-form-item label="变动标题" prop="title" v-if="editForm.pointAction !== 'none'" :rules="editForm.pointAction !== 'none' ? editRules.title : []">
          <el-select v-model="editForm.title" placeholder="请选择或输入变动标题" filterable allow-create default-first-option>
            <el-option label="官方发放" value="官方发放" v-if="editForm.pointAction === 'add'" />
            <el-option label="活动奖励" value="活动奖励" v-if="editForm.pointAction === 'add'" />
            <el-option label="补偿发放" value="补偿发放" v-if="editForm.pointAction === 'add'" />
            <el-option label="系统回收" value="系统回收" v-if="editForm.pointAction === 'subtract'" />
            <el-option label="违规扣除" value="违规扣除" v-if="editForm.pointAction === 'subtract'" />
          </el-select>
        </el-form-item>
        <el-form-item label="修改说明" prop="reason" v-if="editForm.pointAction !== 'none'" :rules="editForm.pointAction !== 'none' ? editRules.reason : []">
          <el-input v-model="editForm.reason" placeholder="请输入积分修改详细说明" />
        </el-form-item>
        <el-form-item label="账号状态" prop="status">
          <el-switch v-model="editForm.status" :active-value="1" :inactive-value="0" active-text="正常" inactive-text="禁用" inline-prompt />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取 消</el-button>
          <el-button type="primary" @click="handleEdit" :loading="editLoading">确认保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 积分记录对话框 -->
    <el-dialog v-model="showPointsRecordDialog" :title="`${currentRecordUser} 的积分记录`" width="700px" destroy-on-close class="custom-dialog">
      <el-table 
        :data="pointsRecordList" 
        v-loading="pointsRecordLoading" 
        style="width: 100%" 
        height="400"
        :header-cell-style="{ background: '#f8fafc', color: '#475569', fontWeight: '600' }"
      >
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="scope">
            <span>{{ new Date(scope.row.created_at).toLocaleString() }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="points_change" label="积分变动" width="120">
          <template #default="scope">
            <span :style="{ color: scope.row.points_change > 0 ? '#67c23a' : '#f56c6c', fontWeight: 'bold' }">
              {{ scope.row.points_change > 0 ? '+' : '' }}{{ parseFloat(Number(scope.row.points_change).toFixed(2)) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="详情" min-width="200">
          <template #default="scope">
            <div style="display: flex; flex-direction: column; gap: 4px;">
              <span style="font-weight: 600; color: #334155;">{{ scope.row.title || '官方发放' }}</span>
              <span style="font-size: 13px; color: #64748b;">{{ scope.row.reason }}</span>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-container" style="margin-top: 15px;">
        <el-pagination
          v-model:current-page="pointsRecordPage"
          v-model:page-size="pointsRecordPageSize"
          :page-sizes="[10, 20, 50]"
          :background="true"
          layout="total, sizes, prev, pager, next"
          :total="pointsRecordTotal"
          @size-change="handlePointsRecordSizeChange"
          @current-change="handlePointsRecordCurrentChange"
        />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, EditPen, Delete, Calendar, User, Search, Message, Phone, Warning, Coin, List } from '@element-plus/icons-vue'

const userList = ref<any[]>([])
const loading = ref(false)
const total = ref(0)

// 搜索和分页
const searchQuery = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

const handleSearch = () => {
  currentPage.value = 1
  fetchUserList()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
  fetchUserList()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchUserList()
}

const getAvatarColor = (name: string) => {
  const colors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#8e44ad', '#16a085', '#d35400'];
  if (!name) return colors[0];
  let hash = 0;
  for (let i = 0; i < name.length; i++) {
    hash = name.charCodeAt(i) + ((hash << 5) - hash);
  }
  return colors[Math.abs(hash) % colors.length];
}

const showAddDialog = ref(false)
const submitLoading = ref(false)
const addFormRef = ref()
const addForm = ref({
  username: '',
  password: '',
  email: '',
  phone: '',
  points: 100.0,
  status: 1
})

const showPointsRecordDialog = ref(false)
const pointsRecordLoading = ref(false)
const pointsRecordList = ref<any[]>([])
const pointsRecordPage = ref(1)
const pointsRecordPageSize = ref(10)
const pointsRecordTotal = ref(0)
const currentRecordUserId = ref<number | null>(null)
const currentRecordUser = ref('')

const fetchPointsRecords = async () => {
  if (!currentRecordUserId.value) return
  pointsRecordLoading.value = true
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/users/points/records?user_id=${currentRecordUserId.value}&page=${pointsRecordPage.value}&page_size=${pointsRecordPageSize.value}`, {
      headers: getHeaders()
    })
    const data = await res.json()
    if (res.ok) {
      pointsRecordList.value = data.list || []
      pointsRecordTotal.value = data.total || 0
    } else {
      ElMessage.error(data.error || '获取积分记录失败')
    }
  } catch (error) {
    ElMessage.error('网络错误，请稍后重试')
  } finally {
    pointsRecordLoading.value = false
  }
}

const openPointsRecordDialog = (row: any) => {
  currentRecordUserId.value = row.id
  currentRecordUser.value = row.username || row.email
  pointsRecordPage.value = 1
  showPointsRecordDialog.value = true
  fetchPointsRecords()
}

const handlePointsRecordSizeChange = (val: number) => {
  pointsRecordPageSize.value = val
  pointsRecordPage.value = 1
  fetchPointsRecords()
}

const handlePointsRecordCurrentChange = (val: number) => {
  pointsRecordPage.value = val
  fetchPointsRecords()
}

const showEditDialog = ref(false)
const editLoading = ref(false)
const editFormRef = ref()
const originalPoints = ref(0)
const editForm = ref({
  id: null as number | null,
  username: '',
  password: '',
  email: '',
  phone: '',
  pointAction: 'none',
  pointsChange: 1,
  title: '官方发放',
  reason: '',
  status: 1
})

const rules = {
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

const editRules = {
  password: [
    { required: false, min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  title: [
    { required: true, message: '请输入或选择变动标题', trigger: ['blur', 'change'] }
  ],
  reason: [
    { required: true, message: '请输入积分修改详细说明', trigger: 'blur' }
  ]
}

const getHeaders = () => {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  }
}

const fetchUserList = async () => {
  loading.value = true
  try {
    const baseUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    const url = new URL(`${baseUrl}/api/admin/users/list`, window.location.origin)
    url.searchParams.append('page', currentPage.value.toString())
    url.searchParams.append('page_size', pageSize.value.toString())
    if (searchQuery.value) {
      url.searchParams.append('search', searchQuery.value)
    }
    
    const res = await fetch(url.toString(), { headers: getHeaders() })
    const data = await res.json()
    if (res.ok) {
      userList.value = data.list || []
      total.value = data.total || 0
    } else {
      ElMessage.error(data.error || '获取用户列表失败')
    }
  } catch (error) {
    ElMessage.error('网络错误，请稍后重试')
  } finally {
    loading.value = false
  }
}

const handleAdd = async () => {
  if (!addFormRef.value) return
  await addFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      submitLoading.value = true
      try {
        const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/users/create`, {
          method: 'POST',
          headers: getHeaders(),
          body: JSON.stringify(addForm.value)
        })
        const data = await res.json()
        if (res.ok) {
          ElMessage.success('添加成功')
          showAddDialog.value = false
          addForm.value = { username: '', password: '', email: '', phone: '', points: 100.0, status: 1 }
          fetchUserList()
        } else {
          ElMessage.error(data.error || '添加失败')
        }
      } catch (error) {
        ElMessage.error('网络错误，请稍后重试')
      } finally {
        submitLoading.value = false
      }
    }
  })
}

const openEditDialog = (row: any) => {
  originalPoints.value = Number(row.points) || 0
  editForm.value = {
    id: row.id,
    username: row.username,
    password: '', // 置空，不填则不改
    email: row.email,
    phone: row.phone,
    pointAction: 'none',
    pointsChange: 1,
    title: '官方发放',
    reason: '',
    status: row.status
  }
  showEditDialog.value = true
}

const handleSubtractClick = (e: MouseEvent) => {
  if (originalPoints.value <= 0) {
    e.preventDefault()
    ElMessage.warning('该用户当前积分为 0，无法扣除')
  }
}

const handlePointActionChange = (val: string) => {
  if (val === 'add') {
    editForm.value.title = '官方发放'
  } else if (val === 'subtract') {
    if (originalPoints.value <= 0) {
      editForm.value.pointAction = 'none'
      return
    }
    editForm.value.title = '系统回收'
    if (editForm.value.pointsChange > originalPoints.value) {
      editForm.value.pointsChange = originalPoints.value
    }
  }
}

const handleEdit = async () => {
  if (!editFormRef.value) return
  await editFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      editLoading.value = true
      try {
        let finalPoints = originalPoints.value
        if (editForm.value.pointAction === 'add') {
          finalPoints += Number(editForm.value.pointsChange)
        } else if (editForm.value.pointAction === 'subtract') {
          const change = Number(editForm.value.pointsChange)
          if (change > originalPoints.value) {
            ElMessage.error(`扣除积分不能超过当前积分 ${originalPoints.value}`)
            editLoading.value = false
            return
          }
          finalPoints -= change
          if (finalPoints < 0) finalPoints = 0
        }

        const payload = {
          ...editForm.value,
          points: finalPoints
        }

        const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/users/update`, {
          method: 'POST',
          headers: getHeaders(),
          body: JSON.stringify(payload)
        })
        const data = await res.json()
        if (res.ok) {
          ElMessage.success('修改成功')
          showEditDialog.value = false
          fetchUserList()
        } else {
          ElMessage.error(data.error || '修改失败')
        }
      } catch (error) {
        ElMessage.error('网络错误，请稍后重试')
      } finally {
        editLoading.value = false
      }
    }
  })
}

const handleDelete = async (row: any) => {
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/users/delete`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ id: row.id })
    })
    const data = await res.json()
    if (res.ok) {
      ElMessage.success('删除成功')
      fetchUserList()
    } else {
      ElMessage.error(data.error || '删除失败')
    }
  } catch (error) {
    ElMessage.error('网络错误，请稍后重试')
  }
}

onMounted(() => {
  fetchUserList()
})
</script>

<style scoped>
.fade-in {
  animation: fadeIn 0.4s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.admin-container {
  padding: 24px;
  background-color: #f1f5f9;
  min-height: calc(100vh - 60px);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  background: #fff;
  padding: 20px 24px;
  border-radius: 12px;
  box-shadow: 0 1px 3px rgba(0,0,0,0.05);
}

.header-info {
  display: flex;
  align-items: center;
  gap: 16px;
}

.icon-wrapper {
  width: 48px;
  height: 48px;
  background: #e6f2ff;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-title h2 {
  margin: 0 0 4px 0;
  font-size: 20px;
  color: #1e293b;
  font-weight: 600;
}

.subtitle {
  font-size: 13px;
  color: #64748b;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.search-input {
  width: 280px;
}

:deep(.search-input .el-input__wrapper) {
  border-radius: 20px;
  box-shadow: 0 0 0 1px #e2e8f0 inset;
}

:deep(.search-input .el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #cbd5e1 inset;
}

:deep(.search-input .el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409eff inset;
}

.add-btn {
  padding: 10px 20px;
  border-radius: 8px;
  font-weight: 500;
  box-shadow: 0 4px 6px -1px rgba(64, 158, 255, 0.2);
  transition: all 0.3s;
}

.add-btn:hover, .add-btn:focus {
  transform: translateY(-1px);
  box-shadow: 0 6px 8px -1px rgba(64, 158, 255, 0.3);
}

.admin-card {
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  background: #fff;
}

.id-badge {
  display: inline-block;
  padding: 2px 8px;
  background: #f1f5f9;
  border-radius: 4px;
  color: #64748b;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 13px;
  font-weight: 500;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-avatar {
  color: #fff;
  font-weight: 600;
  font-size: 16px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: flex-start;
}

.username {
  font-weight: 600;
  color: #334155;
  font-size: 14px;
}

.info-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #475569;
}

.empty-text {
  color: #94a3b8;
}

.status-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 0 10px;
  height: 24px;
  border-radius: 12px;
  border: none;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
}

.dot-success {
  background-color: #409eff;
}

.dot-danger {
  background-color: #ef4444;
}

.time-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  color: #64748b;
  font-size: 13px;
}

.action-btns {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.action-btn {
  width: 32px;
  height: 32px;
  padding: 8px;
  border-radius: 8px;
  background: #f1f5f9;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #e2e8f0;
  transform: scale(1.05);
}

.danger-btn:hover {
  background: #fee2e2;
  color: #ef4444;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  padding: 20px 0 0;
  border-top: 1px solid #f1f5f9;
  margin-top: 20px;
}

:deep(.el-table) {
  --el-table-border-color: #e2e8f0;
  --el-table-row-hover-bg-color: #f8fafc;
}

:deep(.el-table th.el-table__cell) {
  padding: 12px 0;
  border-bottom: 1px solid #e2e8f0;
}

:deep(.el-table td.el-table__cell) {
  padding: 16px 0;
  border-bottom: 1px dashed #e2e8f0;
}

:deep(.custom-table-row) {
  transition: all 0.3s ease;
}

/* 对话框美化 */
.dialog-desc {
  color: #64748b;
  font-size: 14px;
  margin-bottom: 24px;
  margin-top: -10px;
}

:deep(.custom-dialog) {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

:deep(.custom-dialog .el-dialog__header) {
  margin: 0;
  padding: 24px 24px 16px;
}

:deep(.custom-dialog .el-dialog__title) {
  font-weight: 600;
  font-size: 18px;
  color: #1e293b;
}

:deep(.custom-dialog .el-dialog__body) {
  padding: 0 24px 24px;
}

:deep(.custom-dialog .el-dialog__footer) {
  padding: 16px 24px;
  background-color: #f8fafc;
  border-top: 1px solid #e2e8f0;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 0 0 1px #e2e8f0 inset;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #cbd5e1 inset;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #409eff inset !important;
}

:deep(.el-button) {
  border-radius: 8px;
}
</style>
