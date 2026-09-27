<template>
  <div class="inspirations-container">
    <div class="page-header">
      <h2>探索管理</h2>
      <el-button type="primary" @click="showAddDialog">
        <el-icon><Plus /></el-icon> 添加提示词
      </el-button>
    </div>

    <el-card class="box-card" shadow="never">
      <div class="toolbar">
        <div class="search-bar">
          <el-input
            v-model="searchQuery"
            placeholder="搜索提示词内容"
            clearable
            @clear="handleSearch"
            @keyup.enter="handleSearch"
            class="search-input"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button type="danger" :disabled="selectedIds.length === 0" @click="handleBatchDelete" :loading="batchDeleting">批量删除</el-button>
        </div>
      </div>

      <el-table :data="inspirations" style="width: 100%" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="示例图片" width="120">
          <template #default="{ row }">
            <el-image 
              v-if="row.image_url"
              :src="row.image_url" 
              :preview-src-list="[row.image_url]"
              fit="cover" 
              style="width: 60px; height: 60px; border-radius: 4px;" 
              preview-teleported
            />
            <span v-else style="color: #999; font-size: 12px;">无图片</span>
          </template>
        </el-table-column>
        <el-table-column prop="main_category" label="主分类" width="120" />
        <el-table-column prop="sub_category" label="副分类" width="120" />
        <el-table-column prop="content" label="提示词内容" min-width="400">
          <template #default="{ row }">
            <span style="white-space: pre-wrap; word-break: break-all;">{{ row.content }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" size="small" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm
              title="确定删除此提示词吗？"
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

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑提示词' : '添加提示词'" width="500px">
      <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="示例图片" prop="image_url">
          <el-input v-model="form.image_url" placeholder="请输入图片URL（选填）" clearable />
          <div style="margin-top: 10px;" v-if="form.image_url">
            <el-image :src="form.image_url" style="width: 100px; height: 100px; border-radius: 4px;" fit="cover" />
          </div>
        </el-form-item>
        <el-form-item label="主分类" prop="main_category">
          <el-select v-model="form.main_category" filterable allow-create clearable placeholder="请选择或输入主分类" @change="handleMainCategoryChange">
            <el-option v-for="item in mainCategories" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="副分类" prop="sub_category">
          <el-select v-model="form.sub_category" filterable allow-create clearable placeholder="请选择或输入副分类">
            <el-option v-for="item in subCategories" :key="item" :label="item" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="提示词内容" prop="content">
          <el-input 
            v-model="form.content" 
            type="textarea" 
            :rows="6"
            placeholder="请输入提示词内容，这将显示在前台迷你边栏的探索功能中" 
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitForm" :loading="submitting">
            确定
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance } from 'element-plus'
import dayjs from 'dayjs'
import axios from 'axios'

const getAuthHeaders = () => {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

const getApiUrl = (path: string) => {
  return `${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api${path}`
}

interface Inspiration {
  id: number
  content: string
  image_url: string
  main_category: string
  sub_category: string
  created_at: string
  updated_at: string
}

const loading = ref(false)
const inspirations = ref<Inspiration[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchQuery = ref('')
const selectedIds = ref<number[]>([])
const batchDeleting = ref(false)

const mainCategories = ref<string[]>([])
const subCategories = ref<string[]>([])
const categoryTree = ref<any[]>([])

const loadCategories = async () => {
  try {
    const res = await axios.get(getApiUrl('/admin/settings/get?key_name=inspiration_categories'), {
      headers: getAuthHeaders()
    })
    
    if (res.data?.data?.key_value) {
      const parsed = JSON.parse(res.data.data.key_value)
      categoryTree.value = parsed
      mainCategories.value = parsed.map((c: any) => c.name)
    } else {
      // Default fallback
      const defaultMain = ['人物', '风景', '建筑', '科幻', '二次元', '其他']
      mainCategories.value = defaultMain
      categoryTree.value = defaultMain.map(name => ({ name, sub: ['其他'] }))
    }
  } catch (error: any) {
    console.error('获取分类失败', error)
  }
}

const handleMainCategoryChange = (val: string) => {
  form.value.sub_category = ''
  const target = categoryTree.value.find((c: any) => c.name === val)
  if (target && target.sub) {
    subCategories.value = target.sub
  } else {
    subCategories.value = []
  }
}

const dialogVisible = ref(false)
const submitting = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const form = ref({
  id: 0,
  content: '',
  image_url: '',
  main_category: '',
  sub_category: ''
})

const rules = {
  content: [
    { required: true, message: '请输入提示词内容', trigger: 'blur' }
  ]
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

const loadData = async () => {
  loading.value = true
  try {
    const res = await axios.get(getApiUrl('/admin/inspirations/list'), {
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
        search: searchQuery.value
      },
      headers: getAuthHeaders()
    })
    inspirations.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '获取数据失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  loadData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadData()
}

const handleSelectionChange = (selection: Inspiration[]) => {
  selectedIds.value = selection.map(item => item.id)
}

const showAddDialog = () => {
  isEdit.value = false
  form.value = {
    id: 0,
    content: '',
    image_url: '',
    main_category: '',
    sub_category: ''
  }
  dialogVisible.value = true
  if (formRef.value) {
    formRef.value.clearValidate()
  }
}

const handleEdit = (row: Inspiration) => {
  isEdit.value = true
  form.value = {
    id: row.id,
    content: row.content,
    image_url: row.image_url || '',
    main_category: row.main_category || '',
    sub_category: row.sub_category || ''
  }
  
  if (row.main_category) {
    const target = categoryTree.value.find((c: any) => c.name === row.main_category)
    if (target && target.sub) {
      subCategories.value = target.sub
    } else {
      subCategories.value = []
    }
  } else {
    subCategories.value = []
  }
  
  dialogVisible.value = true
}

const submitForm = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value) {
          await axios.post(getApiUrl('/admin/inspirations/update'), {
            id: form.value.id,
            content: form.value.content,
            image_url: form.value.image_url,
            main_category: form.value.main_category,
            sub_category: form.value.sub_category
          }, { headers: getAuthHeaders() })
          ElMessage.success('更新成功')
        } else {
          await axios.post(getApiUrl('/admin/inspirations/create'), {
            content: form.value.content,
            image_url: form.value.image_url,
            main_category: form.value.main_category,
            sub_category: form.value.sub_category
          }, { headers: getAuthHeaders() })
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        loadData()
        loadCategories()
      } catch (error: any) {
        ElMessage.error(error.response?.data?.error || (isEdit.value ? '更新失败' : '添加失败'))
      } finally {
        submitting.value = false
      }
    }
  })
}

const handleDelete = async (id: number) => {
  try {
    await axios.post(getApiUrl('/admin/inspirations/delete'), { id }, { headers: getAuthHeaders() })
    ElMessage.success('删除成功')
    if (inspirations.value.length === 1 && currentPage.value > 1) {
      currentPage.value--
    }
    loadData()
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '删除失败')
  }
}

const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.value.length} 个提示词吗？`, '警告', {
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })
    
    batchDeleting.value = true
    await axios.post(getApiUrl('/admin/inspirations/batch-delete'), { ids: selectedIds.value }, { headers: getAuthHeaders() })
    ElMessage.success('批量删除成功')
    
    // 如果当前页的数据全部被删除，且不是第一页，则跳到前一页
    if (selectedIds.value.length === inspirations.value.length && currentPage.value > 1) {
      currentPage.value--
    }
    loadData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.response?.data?.error || '批量删除失败')
    }
  } finally {
    batchDeleting.value = false
  }
}

onMounted(() => {
  loadData()
  loadCategories()
})
</script>

<style scoped>
.inspirations-container {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  font-size: 20px;
  color: #303133;
}

.box-card {
  border-radius: 8px;
}

.toolbar {
  margin-bottom: 20px;
}

.search-bar {
  display: flex;
  gap: 12px;
}

.search-input {
  width: 300px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
