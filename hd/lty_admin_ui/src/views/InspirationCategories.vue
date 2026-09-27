<template>
  <div class="settings-container fade-in">
    <div class="page-header">
      <div class="header-info">
        <div class="icon-wrapper">
          <el-icon :size="24" color="#409eff"><Folder /></el-icon>
        </div>
        <div class="header-title">
          <h2>探索分类设置</h2>
          <span class="subtitle">配置探索管理中的主分类（一级）和副分类（二级）。</span>
        </div>
      </div>
    </div>

    <el-card shadow="never" class="settings-card">
      <div class="card-header-actions">
        <el-button type="primary" @click="addMainCategory" class="add-main-btn">
          <el-icon><Plus /></el-icon> 添加主分类
        </el-button>
        <el-button type="success" @click="handleSaveCategories" :loading="saving" class="save-btn">
          <el-icon><Check /></el-icon> 保存所有更改
        </el-button>
      </div>

      <div v-loading="loading" class="categories-editor">
        <el-empty v-if="inspirationCategories.length === 0" description="暂无分类数据，请点击上方按钮添加" />
        
        <div v-else class="categories-grid">
          <div v-for="(cat, index) in inspirationCategories" :key="index" class="category-item-card">
            <div class="category-header">
              <div class="header-left">
                <el-icon class="drag-icon"><Grid /></el-icon>
                <el-input 
                  v-model="cat.name" 
                  placeholder="请输入主分类名称" 
                  class="main-category-input"
                  clearable
                >
                  <template #prepend>主分类</template>
                </el-input>
              </div>
              <el-popconfirm
                title="确定要删除该主分类及其所有副分类吗？"
                @confirm="removeMainCategory(index)"
                width="240"
              >
                <template #reference>
                  <el-button type="danger" link class="delete-btn">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </template>
              </el-popconfirm>
            </div>
            
            <div class="category-divider"></div>
            
            <div class="category-body">
              <div class="sub-category-title">副分类标签 ({{ cat.sub.length }})</div>
              <div class="sub-category-list">
                <el-tag
                  v-for="(sub, subIndex) in cat.sub"
                  :key="subIndex"
                  closable
                  effect="light"
                  :disable-transitions="false"
                  @close="removeSubCategory(index, subIndex)"
                  class="sub-tag"
                  type="info"
                >
                  {{ sub }}
                </el-tag>
                <el-input
                  v-if="cat.inputVisible"
                  :ref="(el: any) => setInputRef(el, index)"
                  v-model="cat.inputValue"
                  class="sub-input"
                  size="small"
                  placeholder="输入后回车"
                  @keyup.enter="handleInputConfirm(index)"
                  @blur="handleInputConfirm(index)"
                />
                <el-button v-else class="button-new-tag" size="small" type="primary" plain @click="showInput(index)">
                  <el-icon><Plus /></el-icon> 添加副分类
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { Folder, Delete, Plus, Check, Grid } from '@element-plus/icons-vue'

const loading = ref(false)
const saving = ref(false)

interface CategoryItem {
  name: string
  sub: string[]
  inputVisible?: boolean
  inputValue?: string
}

const inspirationCategories = ref<CategoryItem[]>([])
const inputRefs = ref<any[]>([])

const setInputRef = (el: any, index: number) => {
  if (el) {
    inputRefs.value[index] = el
  }
}

const removeMainCategory = (index: number) => {
  inspirationCategories.value.splice(index, 1)
}

const addMainCategory = () => {
  inspirationCategories.value.push({
    name: '',
    sub: [],
    inputVisible: false,
    inputValue: ''
  })
}

const removeSubCategory = (mainIndex: number, subIndex: number) => {
  inspirationCategories.value[mainIndex].sub.splice(subIndex, 1)
}

const showInput = (index: number) => {
  inspirationCategories.value[index].inputVisible = true
  nextTick(() => {
    if (inputRefs.value[index]) {
      inputRefs.value[index].focus()
    }
  })
}

const handleInputConfirm = (index: number) => {
  const cat = inspirationCategories.value[index]
  if (cat.inputValue) {
    const val = cat.inputValue.trim()
    if (val && !cat.sub.includes(val)) {
      cat.sub.push(val)
    }
  }
  cat.inputVisible = false
  cat.inputValue = ''
}

const getHeaders = () => {
  const token = localStorage.getItem('token')
  return {
    'Authorization': `Bearer ${token}`
  }
}

const fetchSetting = async (keyName: string) => {
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/settings/get?key_name=${keyName}`, {
      headers: { ...getHeaders(), 'Content-Type': 'application/json' }
    })
    if (res.ok) {
      const json = await res.json()
      return json.data?.key_value || ''
    }
  } catch (error) {
    console.error(`Failed to fetch setting ${keyName}`, error)
  }
  return ''
}

const saveSetting = async (keyName: string, keyValue: string, description: string) => {
  try {
    const res = await fetch(`${(window as any).APP_CONFIG?.API_BASE_URL || ''}/api/admin/settings/save`, {
      method: 'POST',
      headers: { ...getHeaders(), 'Content-Type': 'application/json' },
      body: JSON.stringify({
        key_name: keyName,
        key_value: keyValue,
        description: description
      })
    })
    if (!res.ok) {
      throw new Error(`Failed to save ${keyName}`)
    }
  } catch (error) {
    throw error
  }
}

const handleSaveCategories = async () => {
  saving.value = true
  try {
    const validCats = inspirationCategories.value
      .filter(c => c.name.trim())
      .map(c => ({
        name: c.name.trim(),
        sub: c.sub
      }))
    
    await saveSetting('inspiration_categories', JSON.stringify(validCats), '探索分类设置')
    ElMessage.success('分类设置保存成功')
  } catch (error) {
    ElMessage.error('分类设置保存失败')
  } finally {
    saving.value = false
  }
}

const loadSettings = async () => {
  loading.value = true
  try {
    const catsStr = await fetchSetting('inspiration_categories')
    if (catsStr) {
      try {
        const parsed = JSON.parse(catsStr)
        inspirationCategories.value = parsed.map((c: any) => ({
          name: c.name,
          sub: c.sub || [],
          inputVisible: false,
          inputValue: ''
        }))
      } catch (e) {
        console.error('Failed to parse categories', e)
      }
    } else {
        inspirationCategories.value = []
    }
  } catch (error) {
    ElMessage.error('获取配置失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadSettings()
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

.settings-container {
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

.settings-card {
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  background: #fff;
  min-height: 500px;
  padding: 0;
}

.card-header-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-bottom: 1px solid #e2e8f0;
  background-color: #f8fafc;
  border-radius: 12px 12px 0 0;
}

.categories-editor {
  padding: 24px;
}

.categories-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.category-item-card {
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  background: #fff;
  transition: all 0.3s ease;
  overflow: hidden;
}

.category-item-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  border-color: #cbd5e1;
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  background-color: #f8fafc;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.drag-icon {
  color: #94a3b8;
  cursor: grab;
  font-size: 18px;
}

.main-category-input {
  flex: 1;
  max-width: 200px;
}

.main-category-input :deep(.el-input-group__prepend) {
  background-color: #fff;
  color: #64748b;
  font-weight: 500;
  padding: 0 12px;
}

.delete-btn {
  font-size: 16px;
  padding: 8px;
  opacity: 0.6;
  transition: opacity 0.2s;
}

.delete-btn:hover {
  opacity: 1;
}

.category-divider {
  height: 1px;
  background-color: #e2e8f0;
}

.category-body {
  padding: 16px;
}

.sub-category-title {
  font-size: 12px;
  color: #64748b;
  margin-bottom: 12px;
  font-weight: 500;
}

.sub-category-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  min-height: 32px;
}

.sub-tag {
  border-radius: 6px;
  padding: 0 10px;
  height: 28px;
  line-height: 26px;
}

.sub-input {
  width: 120px;
}

.button-new-tag {
  border-radius: 6px;
  border-style: dashed;
}

.save-btn, .add-main-btn {
  border-radius: 8px;
  font-weight: 500;
}
</style>