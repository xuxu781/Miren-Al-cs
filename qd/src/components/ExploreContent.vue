<template>
  <main class="explore-main">
    <!-- 现代科技感光晕背景 -->
    <div class="absolute top-0 left-0 right-0 h-[600px] overflow-hidden pointer-events-none z-0 flex justify-center items-start opacity-70">
      <div class="absolute top-[-50px] left-[15%] w-[40vw] max-w-[500px] h-[300px] rounded-full mix-blend-multiply" style="background-color: rgba(99, 102, 241, 0.3); filter: blur(90px);"></div>
      <div class="absolute top-[20px] right-[15%] w-[35vw] max-w-[450px] h-[300px] rounded-full mix-blend-multiply" style="background-color: rgba(168, 85, 247, 0.2); filter: blur(90px);"></div>
      <div class="absolute top-[80px] left-[30%] w-[40vw] max-w-[500px] h-[250px] rounded-full mix-blend-multiply" style="background-color: rgba(56, 189, 248, 0.25); filter: blur(90px);"></div>
    </div>
    
    <header class="main-header">
      <div class="header-left">
        <div class="mobile-menu-btn" @click="$emit('toggle-menu')">
          <el-icon><Expand /></el-icon>
        </div>
        <div class="header-new-chat-icon desktop-only" @click="$emit('new-chat')" title="新建会话">
          <el-icon><Edit /></el-icon>
        </div>
        <div v-if="!isLoggedIn" class="header-login-btn mobile-only" @click="authStore.openLogin()">
          登录
        </div>
        <div v-else class="header-points-btn mobile-only" @click="$emit('points-click')" title="查看积分记录">
          <el-icon><Coin /></el-icon>
          <span>{{ parseFloat(Number(userPoints || 0).toFixed(2)) }}</span>
        </div>
      </div>
      <div class="header-actions">
        <div class="header-new-chat-icon mobile-only" @click="$emit('new-chat')" title="新建会话">
          <el-icon><Edit /></el-icon>
        </div>
        <div v-if="!isLoggedIn" class="header-login-btn desktop-only" @click="authStore.openLogin()">
          登录 / 注册
        </div>
        <div v-else class="header-points-btn desktop-only" @click="$emit('points-click')" title="查看积分记录">
          <el-icon><Coin /></el-icon>
          <span>{{ parseFloat(Number(userPoints || 0).toFixed(2)) }} 积分</span>
        </div>
      </div>
    </header>

    <div class="explore-content-area">
      <div class="explore-container">
        
        <!-- Hero Section with Input Box -->
        <div class="relative z-50 w-full px-0 sm:px-6 md:px-8 mb-12 mt-0 trae-browser-inspect-draggable">
          <div class="mx-auto w-full max-w-[1454px]">
            <h1 class="relative z-[60] mb-6 md:mb-12 mt-4 md:mt-6 flex flex-col gap-1.5 md:gap-4 text-center font-sans items-center justify-center group">
              <span class="relative block text-[28px] sm:text-[44px] md:text-[64px] font-black leading-[1.2] tracking-tight px-2 flex flex-wrap justify-center items-center">
                <span class="text-[#0f172a]">探索视觉边界</span>
                <span class="text-[#94a3b8] font-light mx-2 md:mx-4 transform -translate-y-[2px]">·</span>
                <span class="text-[#4d6bfe]">Miren AI</span>
              </span>
            </h1>
            <div class="mx-auto w-full max-w-[1200px]">
              <div class="rounded-2xl transition-all duration-200 ease-in-out min-w-0 max-w-full p-0 flex min-h-0 flex-col">
                <div class="relative w-full flex min-h-0 flex-col">
                  <!-- Main Input Box -->
                  <div class="explore-input-wrapper">
                    <!-- 嵌入与 Generate.vue 相同的输入框和工具栏 -->
                    <div class="input-container">
                      <div class="input-tools-container">
                        <div class="input-tools">
                          <el-dropdown class="model-dropdown-wrapper" trigger="click" placement="top-start" :teleported="true" @command="(val: string) => props.form.series_id = val">
                          <button class="combined-settings-btn model-settings-btn" type="button">
                            <div class="model-btn-content">
                              <el-icon class="cpu-icon"><Cpu /></el-icon>
                              <span class="btn-text model-name-text" :title="props.availableModels.find(m => m.series_id === props.form.series_id)?.name || '默认模型'">
                                {{ props.availableModels.find(m => m.series_id === props.form.series_id)?.name || '默认模型' }}
                              </span>
                              <span v-if="props.availableModels.find(m => m.series_id === props.form.series_id)?.activity_tag" class="shimmer-tag" :style="{ fontSize: '10px', color: '#fff', background: props.availableModels.find(m => m.series_id === props.form.series_id)?.activity_tag_color || '#10b981', padding: '2px 4px', borderRadius: '4px', marginLeft: '6px', whiteSpace: 'nowrap', flexShrink: 0 }">
                                {{ props.availableModels.find(m => m.series_id === props.form.series_id)?.activity_tag }}
                              </span>
                              <span v-else-if="props.availableModels.find(m => m.series_id === props.form.series_id) && props.hasFreeResolution(props.availableModels.find(m => m.series_id === props.form.series_id))" class="shimmer-tag" style="font-size: 10px; color: #fff; background: #10b981; padding: 2px 4px; border-radius: 4px; margin-left: 6px; white-space: nowrap; flex-shrink: 0;">
                                限时免费
                              </span>
                            </div>
                            <el-icon class="arrow-icon"><ArrowDown /></el-icon>
                          </button>
                          <template #dropdown>
                            <el-dropdown-menu class="model-dropdown-menu">
                              <el-dropdown-item 
                                v-for="model in props.availableModels" 
                                :key="model.series_id" 
                                :command="model.series_id"
                                :class="{ 'is-active-model': props.form.series_id === model.series_id }"
                              >
                                <div style="display: flex; align-items: center; justify-content: space-between; width: 100%; gap: 12px;">
                                    <span>{{ model.name }}</span>
                                    <span v-if="model.activity_tag" class="shimmer-tag" :style="{ fontSize: '11px', color: '#fff', background: model.activity_tag_color || '#10b981', padding: '2px 6px', borderRadius: '4px', lineHeight: '1.2', whiteSpace: 'nowrap' }">{{ model.activity_tag }}</span>
                                    <span v-else-if="props.hasFreeResolution(model)" class="shimmer-tag" style="font-size: 11px; color: #fff; background: #10b981; padding: 2px 6px; border-radius: 4px; line-height: 1.2; white-space: nowrap;">限时免费</span>
                                  </div>
                              </el-dropdown-item>
                              <el-dropdown-item 
                                v-if="props.availableModels.length === 0" 
                                command="default"
                                :class="{ 'is-active-model': props.form.series_id === 'default' }"
                              >
                                默认模型
                              </el-dropdown-item>
                            </el-dropdown-menu>
                          </template>
                        </el-dropdown>
                        
                        <el-popover placement="top-start" :width="props.windowWidth <= 768 ? props.windowWidth - 32 : 360" trigger="click" popper-class="settings-popover" :popper-options="{ modifiers: [{ name: 'preventOverflow', options: { padding: 16 } }] }">
                          <template #reference>
                            <button class="combined-settings-btn" type="button">
                              <el-icon><Crop /></el-icon>
                              <span class="btn-text">
                                <template v-if="props.currentModelRatios.length > 0">
                                  {{ props.form.aspect_ratio || 'auto' }}
                                  <div class="divider"></div>
                                </template>
                                {{ props.form.resolution || '1K' }}
                                <template v-if="props.currentModelImageCounts.length > 0">
                                  <div class="divider"></div>
                                  <span>{{ props.form.num_images || 1 }}<span class="unit-text">张</span></span>
                                </template>
                              </span>
                            </button>
                          </template>
                          <div class="settings-panel">
                            <div class="setting-item" v-if="props.currentModelRatios.length > 0">
                              <div class="setting-label">比例</div>
                              <div class="setting-options">
                                <div class="option-btn" v-for="ratio in props.currentModelRatios" :key="ratio" :class="{ active: props.form.aspect_ratio === ratio }" @click="props.form.aspect_ratio = ratio">
                                  {{ ratio }}
                                </div>
                              </div>
                            </div>
                            <div class="setting-item">
                              <div class="setting-label">分辨率</div>
                              <div class="setting-options">
                                <div class="option-btn" v-for="res in props.currentModelResolutions" :key="res" :class="{ active: props.form.resolution === res }" @click="props.form.resolution = res">
                                  <div style="display: flex; align-items: center; gap: 6px;">
                                    <span>{{ res }}</span>
                                    <span v-if="props.getTierConfigForRes(props.currentModelInfo, res, 'enabled', true) === 'maintenance'" class="shimmer-tag" style="font-size: 11px; color: #fff; background: #9ca3af; padding: 0 4px; border-radius: 4px; line-height: 1.4; white-space: nowrap;">维护中</span>
                                    <span v-else-if="props.getTierConfigForRes(props.currentModelInfo, res, 'tag', '')" class="shimmer-tag" :style="{ fontSize: '11px', color: '#fff', background: props.getTierConfigForRes(props.currentModelInfo, res, 'tag_color', '#10b981'), padding: '0 4px', borderRadius: '4px', lineHeight: '1.4', whiteSpace: 'nowrap' }">{{ props.getTierConfigForRes(props.currentModelInfo, res, 'tag', '') }}</span>
                                    <span v-else-if="Number(props.getTierConfigForRes(props.currentModelInfo, res, 'credits_per_image', 1)) === 0" class="shimmer-tag" style="font-size: 11px; color: #fff; background: #10b981; padding: 0 4px; border-radius: 4px; line-height: 1.4; white-space: nowrap;">限时免费</span>
                                    <span v-else style="font-size: 11px; opacity: 0.7;">{{ parseFloat(Number(props.getTierConfigForRes(props.currentModelInfo, res, 'credits_per_image', 1)).toFixed(2)) }}积分</span>
                                  </div>
                                </div>
                              </div>
                            </div>
                            <div class="setting-item" v-if="props.currentModelImageCounts.length > 0">
                              <div class="setting-label">生成数量</div>
                              <div class="setting-options">
                                <div class="option-btn" v-for="count in props.currentModelImageCounts" :key="count" :class="{ active: props.form.num_images === Number(count) }" @click="props.form.num_images = Number(count)">
                                  {{ count }}张
                                </div>
                              </div>
                            </div>
                          </div>
                        </el-popover>
                      </div>
                    </div>
                      
                    <div class="input-box" :class="{ 'has-refs': props.form.reference_image && props.form.reference_image.length > 0, 'is-focused': isInputFocused }" @click="isInputShrunk ? (isInputFocused = true, isInputShrunk = false) : null">
                        <div class="reference-upload-item initial-upload-item" @click.stop="isInputShrunk ? (isInputFocused = true, isInputShrunk = false) : props.triggerUpload()" v-if="(!props.form.reference_image || props.form.reference_image.length === 0) && props.currentModelMaxRefImages > 0">
                          <div class="reference-upload-content" style="transform: rotate(8deg);">
                            <svg width="1em" height="1em" viewBox="0 0 24 24" preserveAspectRatio="xMidYMid meet" fill="none" role="presentation" xmlns="http://www.w3.org/2000/svg" class="upload-icon-svg"><g><path data-follow-fill="currentColor" d="M10.8 20a1.2 1.2 0 0 0 2.4 0v-6.8H20a1.2 1.2 0 1 0 0-2.4h-6.8V4a1.2 1.2 0 0 0-2.4 0v6.8H4a1.2 1.2 0 0 0 0 2.4h6.8V20Z" clip-rule="evenodd" fill-rule="evenodd" fill="currentColor"></path></g></svg>
                          </div>
                        </div>
                        
                        <div class="reference-image-preview-inline" v-if="props.form.reference_image && props.form.reference_image.length > 0" :style="{ '--total-items': (props.form.reference_image.length + (props.form.reference_image.length < props.currentModelMaxRefImages ? 1 : 0)) } as any">
                          <div class="preview-list-inline" :class="{ 'is-expandable': props.form.reference_image.length > 1 || (props.form.reference_image.length === 1 && props.currentModelMaxRefImages > 1), 'is-mobile-expanded': isMobileRefExpanded }" @click="isInputShrunk ? (isInputFocused = true, isInputShrunk = false) : (props.windowWidth <= 768 ? isMobileRefExpanded = true : null)">
                              <div class="reference-group-hover-trigger"></div>
                              <div class="preview-container-inline" v-for="(img, index) in props.form.reference_image" :key="index" :style="{ '--index': index, '--rotate': `${index === 0 ? 0 : [6, -4, 2, -8, 8, -6, 4, -2, 10, -10, 5, -5, 7, -7, 3, -3][(Number(index) - 1) % 16]}deg`, zIndex: props.form.reference_image.length - Number(index) + 1 } as any">
                                <el-image 
                                  :src="img" 
                                  class="preview-img-inline" 
                                  :preview-src-list="props.windowWidth > 768 ? props.form.reference_image : []" 
                                  :initial-index="index"
                                  fit="cover"
                                  :preview-teleported="true"
                                  :hide-on-click-modal="true"
                                />
                                <div class="remove-btn-inline" @click.stop="props.clearReferenceImage(Number(index))">
                                  <el-icon><Close /></el-icon>
                                </div>
                                <div class="collapsed-upload-badge" @click.stop="props.triggerUpload" v-if="index === 0 && props.form.reference_image.length < props.currentModelMaxRefImages">
                                  <svg width="1em" height="1em" viewBox="0 0 24 24" preserveAspectRatio="xMidYMid meet" fill="none" role="presentation" xmlns="http://www.w3.org/2000/svg"><g><path data-follow-fill="currentColor" d="M10.8 20a1.2 1.2 0 0 0 2.4 0v-6.8H20a1.2 1.2 0 1 0 0-2.4h-6.8V4a1.2 1.2 0 0 0-2.4 0v6.8H4a1.2 1.2 0 0 0 0 2.4h6.8V20Z" clip-rule="evenodd" fill-rule="evenodd" fill="currentColor"></path></g></svg>
                                </div>
                              </div>
                              <div class="reference-upload-item-inline" @click.stop="props.triggerUpload" v-if="props.form.reference_image.length < props.currentModelMaxRefImages" :style="{ '--index': props.form.reference_image.length, '--rotate': `-12deg`, zIndex: 0 } as any">
                                <div class="reference-upload-content" style="transform: rotate(12deg);">
                                  <svg width="1em" height="1em" viewBox="0 0 24 24" preserveAspectRatio="xMidYMid meet" fill="none" role="presentation" xmlns="http://www.w3.org/2000/svg" class="upload-icon-svg"><g><path data-follow-fill="currentColor" d="M10.8 20a1.2 1.2 0 0 0 2.4 0v-6.8H20a1.2 1.2 0 1 0 0-2.4h-6.8V4a1.2 1.2 0 0 0-2.4 0v6.8H4a1.2 1.2 0 0 0 0 2.4h6.8V20Z" clip-rule="evenodd" fill-rule="evenodd" fill="currentColor"></path></g></svg>
                                </div>
                              </div>
                            </div>
                          </div>
                          
                          <input type="file" ref="fileInputRef" accept="image/*" style="display: none" @change="(e) => props.handleFileUpload(e)" multiple />
                          <el-input
                              ref="inputRef"
                              v-model="props.form.prompt"
                              type="textarea"
                              :rows="props.windowWidth > 768 ? 3 : 1"
                              :autosize="{ minRows: props.windowWidth > 768 ? 3 : 1, maxRows: 6 }"
                              placeholder="请输入你想生成的图片描述..."
                              resize="none"
                              class="chat-input"
                              :input-style="props.windowWidth > 768 ? { minHeight: '72px' } : undefined"
                              :class="{ 'has-references': props.form.reference_image && props.form.reference_image.length > 0, 'no-refs-allowed': props.currentModelMaxRefImages === 0 }"
                              @keydown.enter.prevent="handleSend"
                              @focus="isInputFocused = true; isInputShrunk = false"
                              @blur="isInputFocused = false"
                              @click="isInputFocused = true; isInputShrunk = false"
                            />
                          <el-button
                            v-if="props.windowWidth <= 768 && isMobileRefExpanded"
                            type="info"
                            circle
                            class="send-btn"
                            @click.stop="isMobileRefExpanded = false"
                          >
                            <el-icon><Close /></el-icon>
                          </el-button>
                          <el-button
                            v-else
                            type="primary"
                            circle
                            class="send-btn"
                            :class="{ 'maintenance-btn': props.isCurrentResolutionMaintenance }"
                            :disabled="!props.form.prompt.trim() || props.isCurrentResolutionMaintenance"
                            @click="handleSend"
                            :title="props.isCurrentResolutionMaintenance ? '该模型正在维护中，暂时无法生成' : '开始生成'"
                          >
                            <el-icon v-if="!props.isCurrentResolutionMaintenance"><Position /></el-icon>
                            <span v-else style="font-size: 12px; transform: scale(0.9);">维护中</span>
                          </el-button>
                        </div>
                      </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 分类选择区 (美化版) -->
        <div class="category-filter-section" ref="categorySectionRef" v-if="categories.length > 0">
          <div class="main-categories-wrapper">
            <div class="main-categories">
              <div 
                class="category-tab" 
                :class="{ 'active': currentMainCategory === '' }"
                @click="selectMainCategory('')"
              >
                推荐
              </div>
              <div 
                v-for="cat in categories" 
                :key="cat.name"
                class="category-tab"
                :class="{ 'active': currentMainCategory === cat.name }"
                @click="selectMainCategory(cat.name)"
              >
                {{ cat.name }}
              </div>
            </div>
          </div>
          
          <!-- 优雅展开的副分类 -->
          <transition name="el-fade-in-linear">
            <div class="sub-categories" v-if="currentMainCategory && currentSubCategories.length > 0">
              <div class="sub-categories-inner">
                <div 
                  class="sub-category-tag"
                  :class="{ 'active': currentSubCategory === '' }"
                  @click="selectSubCategory('')"
                >
                  全部
                </div>
                <div 
                  v-for="sub in currentSubCategories" 
                  :key="sub"
                  class="sub-category-tag"
                  :class="{ 'active': currentSubCategory === sub }"
                  @click="selectSubCategory(sub)"
                >
                  {{ sub }}
                </div>
              </div>
            </div>
          </transition>
        </div>

        <!-- 瀑布流布局区 (Prompts Gallery Style) -->
        <div class="explore-masonry" v-loading="loadingInspirations">
          <div 
            v-for="item in publicInspirations" 
            :key="item.id" 
            class="prompt-card"
            @click="useSuggestion(item.content)"
          >
            <el-image 
              v-if="item.image_url" 
              :src="item.image_url" 
              fit="cover" 
              loading="lazy"
              class="card-image"
            >
              <template #placeholder>
                <div class="image-placeholder">
                  <el-icon class="is-loading"><Loading /></el-icon>
                </div>
              </template>
              <template #error>
                <div class="image-error">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
            <div v-else class="explore-no-image">
              <el-icon><Picture /></el-icon>
            </div>
            
            <!-- 悬浮时的遮罩与内容 -->
            <div class="card-overlay">
              <div class="overlay-content">
                <p class="card-prompt-text">{{ item.content }}</p>
                <div class="overlay-actions">
                  <el-button type="primary" size="small" class="try-btn" round>
                    <el-icon style="margin-right: 4px"><MagicStick /></el-icon> 画同款
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <el-empty 
          v-if="publicInspirations.length === 0 && !loadingInspirations" 
          description="暂无内容，敬请期待" 
          :image-size="100" 
        />

        <!-- 加载中指示器 -->
        <div class="loading-more" v-if="loadingInspirations && publicInspirations.length > 0">
          <el-icon class="is-loading"><Loading /></el-icon>
          正在加载更多...
        </div>
        
        <!-- 没有更多数据提示 -->
        <div class="no-more-data" v-if="noMoreData && publicInspirations.length > 0">
          <span class="divider-line"></span>
          <span>已经到底啦</span>
          <span class="divider-line"></span>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Expand, Picture, Loading, Coin, Edit, Position, Cpu, Crop, Close, MagicStick } from '@element-plus/icons-vue'
import request from '@/utils/request'
import { useAuthStore } from '@/store/auth'

const props = defineProps<{
  form: any;
  availableModels: any[];
  currentModelInfo: any;
  currentModelRatios: string[];
  currentModelResolutions: string[];
  currentModelImageCounts: (number | string)[];
  currentModelMaxRefImages: number;
  hasFreeResolution: Function;
  getTierConfigForRes: Function;
  isCurrentResolutionMaintenance: boolean;
  windowWidth: number;
  triggerUpload: () => void;
  clearReferenceImage: (index: number) => void;
  handleFileUpload: (e: Event) => void;
}>()

const emit = defineEmits(['toggle-menu', 'points-click', 'new-chat', 'generate'])

const router = useRouter()
const authStore = useAuthStore()

const isLoggedIn = computed(() => !!authStore.token)
const userPoints = computed(() => {
  const user = authStore.userInfo
  return user ? (user.points || 0) : 0
})

const isInputFocused = ref(false)
const isInputShrunk = ref(false)
const isMobileRefExpanded = ref(false)

const handleSend = () => {
  if (props.form.prompt.trim()) {
    // emit('generate') 
    // 如果想要点击后去生成页面，可以直接跳转
    router.push({ path: '/', query: { q: props.form.prompt } })
  }
}

// Inspirations logic
const publicInspirations = ref<{id: number, content: string, image_url: string}[]>([])
const loadingInspirations = ref(false)
const currentPage = ref(1)
const totalInspirations = ref(0)
const noMoreData = ref(false)

// Categories logic
const categories = ref<any[]>([])
const currentMainCategory = ref('')
const currentSubCategory = ref('')

const currentSubCategories = computed(() => {
  if (!currentMainCategory.value) return []
  const target = categories.value.find(c => c.name === currentMainCategory.value)
  return target ? (target.sub || []) : []
})

const fetchCategories = async () => {
  try {
    const res: any = await request.get('/api/public/inspirations/categories')
    if (res.data?.categories) {
      categories.value = res.data.categories
    } else if (res.data?.data?.categories) {
      categories.value = res.data.data.categories
    }
  } catch (error) {
    console.error('Failed to fetch categories', error)
  }
}

const selectMainCategory = (cat: string) => {
  currentMainCategory.value = cat
  currentSubCategory.value = '' // Reset sub category when main changes
  resetAndFetch()
}

const selectSubCategory = (sub: string) => {
  currentSubCategory.value = sub
  resetAndFetch()
}

const categorySectionRef = ref<HTMLElement | null>(null)

const resetAndFetch = () => {
  currentPage.value = 1
  // publicInspirations.value = [] // 移除此行以避免高度塌陷导致的滚动条跳动
  noMoreData.value = false
  fetchPublicInspirations()
}

const fetchPublicInspirations = async () => {
  if (loadingInspirations.value || noMoreData.value) return
  
  loadingInspirations.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: 20
    }
    if (currentMainCategory.value) params.main_category = currentMainCategory.value
    if (currentSubCategory.value) params.sub_category = currentSubCategory.value

    const res: any = await request.get('/api/public/inspirations/list', { params })
    
    let newList: any[] = []
    
    // 处理后端返回的不同数据格式
    if (res.data?.list) {
      newList = res.data.list
      if (res.data.total !== undefined) {
        totalInspirations.value = res.data.total
      }
    } else if (res.data?.data?.list) {
      newList = res.data.data.list
      if (res.data.data.total !== undefined) {
        totalInspirations.value = res.data.data.total
      }
    } else if (Array.isArray(res.data)) {
      newList = res.data
    } else if (Array.isArray(res.data?.data)) {
      newList = res.data.data
    }
    
    if (newList.length === 0) {
      noMoreData.value = true
      if (currentPage.value === 1) {
        publicInspirations.value = []
      }
    } else {
      if (currentPage.value === 1) {
        publicInspirations.value = newList
      } else {
        publicInspirations.value = [...publicInspirations.value, ...newList]
      }
      
      // Check if we loaded all data
      if (newList.length < 20) {
        noMoreData.value = true
      } else {
        currentPage.value++
      }
    }
  } catch (error) {
    console.error('Failed to fetch inspirations', error)
  } finally {
    loadingInspirations.value = false
  }
}

// Scroll handler for infinite scrolling
const handleScroll = () => {
  const scrollContainer = document.querySelector('.explore-content-area')
  if (!scrollContainer) return
  
  const scrollHeight = scrollContainer.scrollHeight
  const scrollTop = scrollContainer.scrollTop
  const clientHeight = scrollContainer.clientHeight
  
  // If user scrolls within 100px of bottom, fetch more
  if (scrollHeight - scrollTop - clientHeight < 100) {
    fetchPublicInspirations()
  }
}

const useSuggestion = (text: string) => {
  props.form.prompt = text
  router.push({ path: '/', query: { q: text } })
}

onMounted(() => {
  fetchCategories()
  fetchPublicInspirations()
  const scrollContainer = document.querySelector('.explore-content-area')
  if (scrollContainer) {
    scrollContainer.addEventListener('scroll', handleScroll)
  }
})

onUnmounted(() => {
  const scrollContainer = document.querySelector('.explore-content-area')
  if (scrollContainer) {
    scrollContainer.removeEventListener('scroll', handleScroll)
  }
})
</script>

<style scoped>
.explore-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;
  background-color: #f9fafb; /* 浅灰底色，贴近豆包 */
  position: relative;
}

/* 样式部分保持 Generate.vue 的 input 样式结构 */
.explore-input-wrapper {
  position: relative;
  width: 100%;
  max-width: 1400px; /* 从 1200px 增加到 1400px，减少左右留白 */
  margin: 0 auto;
  background-color: transparent;
  border-radius: 20px;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px 0 10px 0;
  pointer-events: auto; /* 允许子元素接收点击事件 */
}

.explore-input-wrapper .input-container,
.explore-input-wrapper .input-tools-container {
  pointer-events: auto;
}

.input-area-wrapper.is-shrunk {
  padding-bottom: calc(32px + env(safe-area-inset-bottom)) !important; /* 增加悬浮高度 */
}

.input-area-wrapper.is-shrunk::before {
  opacity: 0;
}

.input-area-wrapper.is-shrunk .input-tools-container {
  height: 0;
  margin: 0;
  opacity: 0;
  overflow: hidden;
  pointer-events: none;
  transform: translateY(10px);
}

.input-area-wrapper.is-shrunk .input-box {
  width: 100%;
  min-height: 40px;
  border-radius: 0;
  box-shadow: none;
  border-color: transparent;
  background-color: transparent;
  backdrop-filter: none;
  padding: 0;
  gap: 8px !important; /* 收缩状态下减小图片与文字之间的间距 */
}

.input-area-wrapper.is-shrunk .initial-upload-item {
  width: 28px !important;
  height: 32px !important;
  margin-right: 0 !important; /* 取消额外边距，完全靠 gap 控制 */
  transform: rotate(-8deg) !important;
  border-radius: 6px;
  position: relative !important;
  pointer-events: auto; /* 允许交互 */
}

.input-area-wrapper.is-shrunk .initial-upload-item:hover,
.input-area-wrapper.is-shrunk .initial-upload-item:active {
  transform: translateY(-4px) scale(1.1) rotate(0deg) !important; /* 相对定位使用 px 向上浮动 */
  z-index: 30 !important;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15) !important;
}

.input-area-wrapper.is-shrunk .initial-upload-item .upload-icon-svg,
.input-area-wrapper.is-shrunk .reference-upload-item-inline .upload-icon-svg {
  font-size: 14px !important; /* 调整图标尺寸 */
}

.input-area-wrapper.is-shrunk .reference-image-preview-inline {
  margin-right: 0 !important; /* 取消额外边距，完全靠 gap 控制 */
  pointer-events: auto; /* 允许交互 */
  overflow: visible !important; /* 收缩状态下强制取消所有滚动条 */
}

.input-area-wrapper.is-shrunk .preview-list-inline {
  width: 28px !important;
  height: 32px !important; /* 恢复为收缩状态应有的高度，因为不展开了 */
  transition: all 0.3s ease;
  pointer-events: auto;
  cursor: pointer;
}

.input-area-wrapper.is-shrunk .preview-container-inline {
  width: 28px !important;
  height: 32px !important;
  border-radius: 6px;
  transform: translateY(-50%) rotate(var(--rotate)) !important;
  pointer-events: none !important; /* 收缩状态下禁止图片本身的点击事件（防全屏预览） */
}

.input-area-wrapper.is-shrunk .reference-upload-item-inline {
  width: 28px !important;
  height: 32px !important;
  border-radius: 6px;
  transform: translateY(-50%) rotate(var(--rotate)) !important;
  pointer-events: none !important;
}

/* 强制收缩状态下的所有图片永远重叠，禁止任何形式的展开和单独悬浮 */
.input-area-wrapper.is-shrunk .preview-list-inline:hover,
.input-area-wrapper.is-shrunk .preview-list-inline.is-expandable:hover,
.input-area-wrapper.is-shrunk .preview-list-inline.is-mobile-expanded,
.input-area-wrapper.is-shrunk .preview-list-inline.is-expandable.is-mobile-expanded {
  width: 28px !important; /* 强制保持收缩状态宽度，禁止展开 */
  margin-right: 0 !important; /* 强制取消 hover/展开 状态下的右侧间距，防止文字位移 */
}

.input-area-wrapper.is-shrunk .preview-list-inline:hover .preview-container-inline,
.input-area-wrapper.is-shrunk .preview-list-inline.is-mobile-expanded .preview-container-inline,
.input-area-wrapper.is-shrunk .preview-list-inline.is-expandable:hover .preview-container-inline,
.input-area-wrapper.is-shrunk .preview-list-inline.is-expandable.is-mobile-expanded .preview-container-inline {
  left: 0 !important; /* 强制所有图片重叠在一起 */
  transform: translateY(-50%) rotate(var(--rotate)) !important; /* 强制取消所有放大和上浮，只保留堆叠 */
  z-index: calc(20 - var(--index)) !important;
}

.input-area-wrapper.is-shrunk .preview-list-inline:hover .reference-upload-item-inline,
.input-area-wrapper.is-shrunk .preview-list-inline.is-mobile-expanded .reference-upload-item-inline,
.input-area-wrapper.is-shrunk .preview-list-inline.is-expandable:hover .reference-upload-item-inline,
.input-area-wrapper.is-shrunk .preview-list-inline.is-expandable.is-mobile-expanded .reference-upload-item-inline {
  display: none !important; /* 收缩状态下隐藏额外上传按钮 */
}

/* 强制取消单张图片的悬浮/点击放大效果 */
.input-area-wrapper.is-shrunk .preview-container-inline:hover,
.input-area-wrapper.is-shrunk .preview-container-inline:active {
  transform: translateY(-50%) rotate(var(--rotate)) !important;
  box-shadow: none !important;
}

.input-area-wrapper.is-shrunk .collapsed-upload-badge,
.input-area-wrapper.is-shrunk .remove-btn-inline {
  display: none !important; /* 收缩状态下隐藏图片上的删除和添加小按钮，保持图标纯净 */
}

.input-area-wrapper.is-shrunk :deep(.chat-input .el-textarea__inner) {
  padding-left: 12px !important;
  font-size: 14px;
  min-height: 24px !important;
  line-height: 24px !important;
}

.input-area-wrapper.is-shrunk .send-btn {
  transform: scale(0.85);
  margin-bottom: 0;
  margin-right: -4px;
}

@media (max-width: 768px) {
  .input-area-wrapper {
    padding-bottom: calc(24px + env(safe-area-inset-bottom));
  }
  .input-area-wrapper.is-shrunk {
    padding-bottom: calc(24px + env(safe-area-inset-bottom)) !important; /* 手机端也适当调高 */
  }

  .content-area {
    padding-bottom: 140px !important; /* 手机端大幅增加底部内边距，确保不被多行输入框遮挡 */
  }
  .result-screen {
    padding-bottom: 40px !important; /* 手机端恢复结果区域底部内边距 */
  }

  .input-area-wrapper.is-shrunk .input-box {
    width: 100%;
    min-width: unset;
  }
  
  .explore-input-wrapper {
    padding: 0 2px !important; /* 手机端极致压缩外层容器的内边距 */
    width: 125%;
    margin-left: -12.5%;
    transform: scale(0.75);
    transform-origin: center top;
    margin-bottom: -20px; /* 补偿缩放后的高度空隙 */
  }

  .input-container {
    padding: 8px 6px 4px !important; /* 手机端进一步压缩底部内边距 */
    border-radius: 16px !important;
  }
  
  .reference-upload-item.initial-upload-item {
    width: 40px;
    height: 48px;
    margin-right: 12px; /* 适当增加初始状态下旋转卡片与文字的间距，防止遮挡 */
  }
  
  .preview-list-inline {
    width: 48px;
    height: 56px; /* 减小高度，避免撑高整个输入框 */
  }
  
  .preview-container-inline {
    width: 48px; /* 恢复为电脑端的宽度 */
    height: 56px; /* 恢复为电脑端的高度 */
    transform: translateY(-50%) rotate(var(--rotate)) !important; /* 恢复原有的旋转效果 */
  }
  
  .input-box .preview-list-inline.is-expandable.is-mobile-expanded {
    width: calc(48px + (var(--total-items, 1) - 1) * 56px); /* 与电脑端间距保持完全一致 */
    margin-right: 16px; /* 恢复适当间距，防止旋转图片遮挡右侧输入框文字 */
  }
  
  .input-box .preview-list-inline.is-expandable.is-mobile-expanded .preview-container-inline {
    left: calc(var(--index) * 56px); /* 与电脑端间距保持完全一致 */
    transform: translateY(-50%) rotate(var(--rotate)) scale(1) !important;
    z-index: calc(20 - var(--index)) !important;
  }
  
  .input-box .preview-list-inline.is-expandable.is-mobile-expanded .preview-container-inline:active {
    transform: translateY(-65%) scale(1.1) rotate(var(--rotate)) !important;
    z-index: 30 !important;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15) !important; /* 增加阴影增强浮动感 */
  }
  
  .input-box .preview-list-inline.is-expandable.is-mobile-expanded .reference-upload-item-inline {
    left: calc(var(--index) * 56px); /* 与电脑端间距保持完全一致 */
    width: 48px; /* 恢复为电脑端的宽度 */
    height: 56px; /* 恢复为电脑端的高度 */
    transform: translateY(-50%) rotate(var(--rotate, 0deg)) !important;
    opacity: 1;
    z-index: 10 !important;
  }
  
  .input-box .preview-list-inline.is-expandable.is-mobile-expanded .reference-upload-item-inline:active {
    background-color: #f3f4f6;
    border-color: #d1d5db;
    color: #4b5563;
    transform: translateY(-65%) scale(1.1) rotate(var(--rotate, 0deg)) !important;
    z-index: 30 !important;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15) !important;
  }
  
  /* 移动端展开时显示删除按钮和隐藏小上传按钮 */
  .preview-list-inline.is-mobile-expanded .remove-btn-inline {
    opacity: 1;
  }
  
  .input-box .preview-list-inline.is-expandable.is-mobile-expanded .collapsed-upload-badge {
    opacity: 0;
    pointer-events: none;
  }

  :deep(.chat-input) {
    min-height: 40px !important;
  }
  
  :deep(.chat-input .el-textarea__inner) {
    min-height: 40px !important;
    padding-top: 8px !important;
    padding-bottom: 8px !important; /* 手机端减小文本框的上下内边距 */
  }

  .input-box {
    min-height: 40px !important; /* 手机端取消强制的 72px 高度限制 */
    padding: 4px 0px !important;
  }

  .send-btn {
    margin-bottom: 0 !important; /* 手机端取消底部外边距 */
  }
}

.input-area-wrapper::before {
  display: none;
}

.reference-image-preview {
  padding: 12px 12px 0 12px;
}

.preview-list {
  position: relative;
  height: 72px;
  padding: 4px 12px;
  background-color: transparent;
  display: flex;
  align-items: center;
  pointer-events: none;
}
.preview-container {
  position: absolute;
  width: 48px;
  height: 64px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e5e7eb;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: center center;
  background-color: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  pointer-events: auto;
}
.preview-container:hover {
  transform: translateY(-8px) rotate(0deg) scale(1.05) !important;
  z-index: 20 !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.preview-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.remove-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 20px;
  height: 20px;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 12px;
  pointer-events: auto;
}

.remove-btn:hover {
  background: rgba(0, 0, 0, 0.7);
}

.reference-upload-item {
  position: absolute;
  width: 48px;
  height: 64px;
  border-radius: 8px;
  background-color: #f9fafb;
  border: 1px dashed #d1d5db;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), box-shadow 0.3s cubic-bezier(0.4, 0, 0.2, 1), background-color 0.3s;
  transform-origin: center center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  pointer-events: auto;
}
.reference-upload-item.initial-upload-item {
  position: relative;
  transform: rotate(-8deg);
  margin-right: 8px;
  width: 48px;
  height: 56px;
  border-radius: 8px;
  border: 1px dashed #e5e7eb;
  background-color: #f9fafb; /* 恢复浅灰色背景 */
  color: #9ca3af;
  flex-shrink: 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04); /* 恢复轻微阴影 */
  z-index: 10;
  pointer-events: auto;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  transform-origin: center center;
  transition: all 0.2s ease;
}
.reference-upload-item.initial-upload-item:hover {
  background-color: #f3f4f6; /* 悬浮时加深背景 */
  border-color: #d1d5db;
  color: #4b5563;
  transform: rotate(0deg) scale(1.02) !important;
  z-index: 20 !important;
}
.reference-upload-item.initial-upload-item .upload-icon-svg {
  font-size: 20px;
}

.reference-upload-item:hover {
  background-color: #f3f4f6;
  border-color: #9ca3af;
  transform: translateY(-8px) rotate(0deg) scale(1.05) !important;
  z-index: 20 !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.reference-upload-content {
  display: flex;
  justify-content: center;
  align-items: center;
}
.upload-icon-svg {
  font-size: 20px;
  color: #6b7280;
}
.initial-upload-item .upload-icon-svg {
  font-size: 20px;
}

.reference-image-preview-inline {
  position: relative;
  z-index: 10;
  display: flex;
  align-items: center;
  margin-right: 4px; /* 电脑端默认减小右侧间距，拉近和文字距离 */
  flex-shrink: 0;
}

.preview-list-inline {
  position: relative;
  width: 48px;
  height: 72px;
  background-color: transparent;
  display: flex;
  align-items: center;
  transition: width 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  overflow: visible;
}

.input-box .preview-list-inline.is-expandable.is-mobile-expanded {
  width: calc(48px + (var(--total-items, 1) - 1) * 56px);
}

.reference-group-hover-trigger {
  position: absolute;
  left: 0;
  top: -10px;
  bottom: -10px;
  width: 100%;
  z-index: -1;
  transition: width 0.3s;
}

.preview-list-inline.is-mobile-expanded .reference-group-hover-trigger {
  width: calc(48px + (var(--total-items, 1) - 1) * 56px);
}

.preview-container-inline {
  position: absolute;
  left: 0;
  top: 50%;
  width: 48px;
  height: 56px;
  border-radius: 8px;
  overflow: hidden;
  border: none;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform: translateY(-50%) rotate(var(--rotate));
  transform-origin: center center;
  background-color: transparent;
  box-shadow: none;
  pointer-events: auto;
}

.input-box .preview-list-inline.is-expandable.is-mobile-expanded .preview-container-inline {
  left: calc(var(--index) * 56px);
  z-index: calc(20 - var(--index)) !important;
}

.preview-img-inline {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid rgba(0, 0, 0, 0.04);
}

.remove-btn-inline {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 20px;
  height: 20px;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 12px;
  pointer-events: auto;
  opacity: 0;
  transition: opacity 0.2s, background 0.2s;
}

.preview-list-inline.is-mobile-expanded .remove-btn-inline {
  opacity: 1;
}

.collapsed-upload-badge {
  position: absolute;
  bottom: 2px;
  right: 2px;
  width: 16px;
  height: 16px;
  background-color: rgba(0, 0, 0, 0.6);
  color: white;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  font-size: 10px;
  pointer-events: auto;
  opacity: 1;
  transition: opacity 0.2s, background-color 0.2s;
  z-index: 10;
}

.input-box .preview-list-inline.is-expandable.is-mobile-expanded .collapsed-upload-badge {
  opacity: 0;
  pointer-events: none;
}

.reference-upload-item-inline {
  position: absolute;
  left: 0;
  top: 50%;
  width: 48px;
  height: 56px;
  border-radius: 8px;
  background-color: #f9fafb;
  border: 1px dashed #e5e7eb;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: center center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  pointer-events: auto;
  opacity: 1;
  color: #9ca3af;
  transform: translateY(-50%) rotate(var(--rotate, 0deg));
}

.input-box .preview-list-inline.is-expandable.is-mobile-expanded .reference-upload-item-inline {
  left: calc(var(--index) * 56px);
  opacity: 1;
  z-index: 10 !important;
}

/* Desktop Hover Expansion Rules */
@media (hover: hover) and (pointer: fine) {
  .input-box .preview-list-inline.is-expandable:hover {
    width: calc(48px + (var(--total-items, 1) - 1) * 56px);
    margin-right: 4px; /* 电脑端悬浮展开时也减小右侧间距 */
  }
  
  .preview-list-inline:hover .reference-group-hover-trigger {
    width: calc(48px + (var(--total-items, 1) - 1) * 56px);
  }

  .input-box .preview-list-inline.is-expandable:hover .preview-container-inline {
    left: calc(var(--index) * 56px);
    z-index: calc(20 - var(--index)) !important;
  }

  .preview-container-inline:hover {
    transform: translateY(-65%) scale(1.1) rotate(var(--rotate)) !important;
    z-index: 30 !important;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15) !important;
  }

  .preview-container-inline:hover .remove-btn-inline {
    opacity: 1;
  }
  
  .remove-btn-inline:hover {
    background: rgba(0, 0, 0, 0.7);
  }

  .collapsed-upload-badge:hover {
    background-color: rgba(0, 0, 0, 0.8);
  }

  .input-box .preview-list-inline.is-expandable:hover .collapsed-upload-badge {
    opacity: 0;
    pointer-events: none;
  }

  .input-box .preview-list-inline.is-expandable:hover .reference-upload-item-inline {
    left: calc(var(--index) * 56px);
    opacity: 1;
    z-index: 10 !important;
  }

  .reference-upload-item-inline:hover {
    background-color: #f3f4f6;
    border-color: #d1d5db;
    color: #4b5563;
    transform: translateY(-65%) scale(1.1) rotate(var(--rotate, 0deg)) !important;
    z-index: 30 !important;
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15) !important;
  }
}



.input-container {
  width: 100%;
  max-width: 1400px; /* 从 1200px 增加到 1400px */
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  border-radius: 24px;
  box-sizing: border-box;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.04), 0 4px 10px rgba(0, 0, 0, 0.02);
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.search-suggestions {
  width: 100%;
  max-width: 1400px;
  margin: 16px auto 0;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  animation: fadeIn 0.5s ease-out;
}

.suggestion-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  justify-content: center;
}

.suggestion-label {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
  margin-right: 4px;
}

.suggestion-tag {
  display: inline-flex;
  align-items: center;
  font-size: 13px;
  color: #4b5563;
  background-color: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(8px);
  padding: 6px 14px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s ease;
  border: 1px solid rgba(229, 231, 235, 0.8);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.02);
}

.suggestion-tag:hover {
  background-color: #ffffff;
  color: #111827;
  border-color: #d1d5db;
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.05);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 768px) {
  .search-suggestions {
    padding: 0 8px;
    margin-top: 24px;
  }
  .suggestion-tags {
    justify-content: flex-start;
    overflow-x: auto;
    flex-wrap: nowrap;
    padding-bottom: 8px;
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
  .suggestion-tags::-webkit-scrollbar {
    display: none;
  }
  .suggestion-tag {
    white-space: nowrap;
    flex-shrink: 0;
  }
  .suggestion-label {
    white-space: nowrap;
    flex-shrink: 0;
  }
}

.input-area-wrapper.is-shrunk .input-container {
  width: 50%;
  min-width: 280px;
  max-width: 500px;
  padding: 8px 16px;
  border-radius: 36px;
  background-color: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
  border-color: transparent;
  gap: 0;
}

@media (max-width: 768px) {
  .input-area-wrapper.is-shrunk .input-container {
    width: 90%;
    min-width: unset;
  }
}

.input-container:focus-within {
  border-color: #d1d5db;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.08), 0 4px 12px rgba(0, 0, 0, 0.04);
}

.model-dropdown-wrapper {
  max-width: 100%;
  display: inline-flex;
  flex: 0 0 auto;
  min-width: 0;
}

.model-settings-btn {
  justify-content: flex-start;
}

.model-btn-content {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  flex: 0 1 auto;
  min-width: 0;
}

.cpu-icon {
  flex-shrink: 0;
}

.model-name-text {
  overflow: visible;
  white-space: nowrap;
  flex: 0 1 auto;
  min-width: 0;
  text-align: left;
  max-width: none; /* 移除最大宽度限制，防止文字被截断 */
  line-height: 1.2;
}

.arrow-icon {
  font-size: 12px !important;
  color: #9ca3af !important;
}

.desktop-only {
  display: flex !important;
}

.mobile-only {
  display: none !important;
}

.input-tools-container {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: bottom center;
  width: 100%;
  max-width: 1400px;
  margin-bottom: 8px;
  display: flex;
  justify-content: flex-start;
  position: relative;
  z-index: 100;
}

.input-tools {
    display: flex;
    gap: 8px;
    padding: 0;
    align-items: center;
    flex-wrap: wrap; /* 允许在空间不足时换行 */
    overflow: visible; /* 移除任何截断 */
  }

:deep(.input-tools .el-select__wrapper) {
  box-shadow: none !important;
  background-color: #f3f4f6;
  border-radius: 8px;
  transition: background-color 0.2s;
}

:deep(.input-tools .el-select__wrapper:hover) {
  background-color: #e5e7eb;
}

:deep(.input-tools .el-select__wrapper.is-focused) {
  background-color: #e5e7eb;
}

.ref-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  border: 1px solid transparent;
  background-color: #f4f4f5;
  border-radius: 20px;
  color: #3f3f46;
  font-weight: 500;
  padding: 4px 14px;
  height: 34px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.ref-btn:hover {
  background-color: #e4e4e7;
  color: #18181b;
  transform: translateY(-1px);
}

.ref-btn:active {
  transform: translateY(0);
}

.ref-btn .el-icon {
  font-size: 15px;
  color: #71717a;
  transition: color 0.25s ease;
}

.ref-btn:hover .el-icon {
  color: #18181b;
}

.size-select {
  width: 120px;
}

.input-box {
  position: relative;
  display: flex;
  align-items: center;
  width: 100%;
  box-sizing: border-box;
  max-width: 1400px; /* 从 1200px 增加到 1400px */
  background: transparent;
  border: none;
  border-radius: 20px;
  padding: 8px 0px; /* 取消左右内边距 */
  gap: 8px; /* 减小内部元素间距 */
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  z-index: 1;
  min-height: 72px;
  box-shadow: none;
}

.input-box:focus-within {
  border-color: transparent;
  box-shadow: none;
}

.input-box.has-refs {
  min-height: 76px;
  padding-bottom: 8px;
}

.input-area-wrapper.is-shrunk .input-box.has-refs {
  min-height: 40px;
  padding-bottom: 0px !important;
}

:deep(.chat-input) {
  flex: 1;
  position: relative;
  z-index: 1;
  background-color: transparent !important;
  --el-input-focus-border-color: transparent;
  --el-input-hover-border-color: transparent;
  --el-input-border-color: transparent;
  min-height: 72px;
}

:deep(.chat-input .el-textarea__inner) {
  padding-left: 8px !important;
  padding-right: 8px !important;
  box-shadow: none !important;
  background-color: transparent !important;
  padding-top: 12px !important;
  padding-bottom: 12px !important;
  font-size: 16px;
  line-height: 24px !important;
  min-height: 72px !important;
  color: #111827;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: none !important;
  display: flex;
  align-items: center;
}

/* 移除了其他的 padding-left 规则 */

:deep(.chat-input .el-textarea__inner::placeholder) {
  color: #9ca3af;
}

.send-btn {
  margin-bottom: 4px;
  background-color: #111827 !important;
  border-color: #111827 !important;
  transition: transform 0.2s ease, opacity 0.2s, width 0.3s ease;
}

.send-btn.maintenance-btn {
  width: auto;
  min-width: 52px;
  border-radius: 16px;
  background-color: #9ca3af !important;
  border-color: #9ca3af !important;
  color: #fff !important;
  opacity: 1 !important;
}

.send-btn:hover:not(.is-disabled) {
  transform: scale(1.05);
  opacity: 0.9;
}

.send-btn.is-disabled {
  background-color: #e5e7eb !important;
  border-color: #e5e7eb !important;
  color: #9ca3af !important;
}

/* 提示框样式优化 */
:deep(.el-tooltip__popper) {
  padding: 8px 12px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
/* Responsive */
@keyframes text-shimmer {
  0% { background-position: 0% 50%; }
  100% { background-position: 100% 50%; }
}

.animate-text-shimmer {
  animation: text-shimmer 3s ease-in-out infinite alternate;
}

.gradient-text {
  background: linear-gradient(110deg, #4F46E5 0%, #9333EA 45%, #EC4899 100%);
  background-size: 250% 100%;
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  color: transparent;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.1));
}

@media (min-width: 769px) {
  .mobile-menu-btn {
    display: none;
  }
  .mobile-only {
    display: none !important;
  }
}

@media (max-width: 768px) {
  .settings-popover {
    width: calc(100vw - 32px) !important;
  }
  .header-actions {
    display: flex;
  }

  .mini-sidebar {
      width: 60px;
      position: absolute;
      top: 0;
      bottom: 0;
      left: 0;
      border-radius: 0;
      transform: translateX(-100%);
      transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
      padding-bottom: calc(48px + env(safe-area-inset-bottom));
      z-index: 1010;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      box-sizing: border-box;
    }
    
    .history-sidebar {
      width: 240px;
      position: absolute;
      top: 0;
      bottom: 0;
      left: 60px;
      border-radius: 0 24px 24px 0;
      transform: translateX(calc(-100% - 60px));
      transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
      padding-bottom: calc(48px + env(safe-area-inset-bottom));
      box-shadow: 8px 0 24px rgba(0, 0, 0, 0.08);
      z-index: 1005;
      box-sizing: border-box;
    }
  
  .mini-sidebar.mobile-open,
  .history-sidebar.mobile-open {
    transform: translateX(0);
    transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
  }
  
  .sidebar-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 1000;
    opacity: 0;
    animation: fadeInOverlay 0.25s ease forwards;
    -webkit-tap-highlight-color: transparent;
  }
  
  @keyframes fadeInOverlay {
  to {
    opacity: 1;
  }
}


.global-account-container {
    width: 300px;
    transform: translateX(-100%);
    transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
    z-index: 2000; /* 确保移动端下也是最高层级 */
  }

  .global-account-container.explore-mode {
    width: 60px;
  }
  
  .global-account-container.mobile-open {
    transform: translateX(0);
    transition: transform 0.3s cubic-bezier(0.4, 0.0, 0.2, 1);
  }
  
  .accountTriggerAvatar {
    width: 60px;
  }

  .main-header {
    padding: 0 16px;
  }

  .mobile-menu-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 20px;
    cursor: pointer;
    color: #4b5563;
    padding: 6px;
    border-radius: 8px;
    background-color: #ffffff;
    border: 1px solid #e5e7eb;
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
    transition: all 0.2s ease;
    -webkit-tap-highlight-color: transparent; /* 消除移动端点击时的默认高亮（蓝色背景） */
    outline: none; /* 消除可能存在的焦点轮廓 */
  }
  
  .mobile-menu-btn:hover {
    background-color: #f9fafb;
    border-color: #d1d5db;
    color: #111827;
  }
  
  .input-tools-container {
    max-width: 100%;
    width: 100%;
    padding: 0;
    box-sizing: border-box;
    margin-bottom: 0;
  }

  .input-container {
    max-width: 100%;
    width: 100%;
    padding: 8px 6px 4px !important; /* 手机端减小顶部和底部 padding，特别是底部 */
    border-radius: 16px;
  }

  .input-area-wrapper {
    padding: 6px 0 6px 0; /* 减小 wrapper 留白 */
  }

  .message-inner {
    max-width: 100%;
  }

  .header-title {
    font-size: 14px;
    max-width: 160px;
  }

  .suggestion-cards {
    flex-direction: column;
    align-items: stretch;
  }

  .desktop-only {
    display: none !important;
  }

  .mobile-only {
    display: flex !important;
  }

  /* 移动端保留花哨效果，并在展开后允许横向滑动 */
  .reference-image-preview-inline {
    position: relative;
    z-index: 10;
    display: flex;
    align-items: center;
    margin-right: 0px; /* 移动端下移除右侧间距，让它更靠近文本 */
    flex-shrink: 0;
    max-width: calc(100vw - 120px); /* 留出输入框和发送按钮的空间 */
    overflow-x: auto;
    overflow-y: visible; /* 允许上下内容溢出，防止截断 */
    scrollbar-width: none;
    -ms-overflow-style: none;
  }

  .reference-image-preview-inline::-webkit-scrollbar {
    display: none !important;
  }

  .reference-group-hover-trigger {
    display: none;
  }

  .input-tools {
    display: flex;
    gap: 8px;
    padding: 0; 
    align-items: center;
    flex-wrap: nowrap;
    overflow-x: auto;
    overflow-y: visible; /* 允许下拉菜单溢出显示 */
    scrollbar-width: none !important;
    -ms-overflow-style: none;
    -webkit-overflow-scrolling: touch;
    margin-bottom: 8px;
    padding-bottom: 4px; /* 为可能出现的轻微溢出留空间 */
  }
  
  /* Webkit 浏览器 (Chrome, Safari, iOS 等) 强制隐藏滚动条 */
  .input-tools::-webkit-scrollbar {
    display: none !important;
    width: 0 !important;
    height: 0 !important;
    background: transparent !important;
  }
}

/* Points Record Dialog Styles */
:deep(.el-overlay-dialog) {
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 彻底重置 Element Plus 弹窗原生样式 */
.points-record-dialog.el-dialog,
:deep(.points-record-dialog.el-dialog) {
  background: transparent !important;
  box-shadow: none !important;
  border-radius: 0 !important;
  padding: 0 !important;
  margin: 0 !important;
  border: none !important;
  --el-dialog-bg-color: transparent !important;
  --el-dialog-box-shadow: none !important;
  --el-dialog-padding-primary: 0 !important;
}

:deep(.points-record-dialog .el-dialog__header) {
  display: none !important;
  padding: 0 !important;
  margin: 0 !important;
}

:deep(.points-record-dialog .el-dialog__body) {
  padding: 0 !important;
  margin: 0 !important;
  background: transparent !important;
}

.points-container {
  background: #ffffff;
  border-radius: 24px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  height: 600px;
  max-height: 80vh;
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.1);
}

.cdkey-redemption-card {
  margin: 0 32px 24px;
}
.cdkey-redemption-card :deep(.el-input-group__append) {
  background-color: var(--theme-color, #FE2C55);
  border-color: var(--theme-color, #FE2C55);
  color: white;
  border-radius: 0 12px 12px 0;
  padding: 0;
  overflow: hidden;
}
.cdkey-redemption-card :deep(.el-input__wrapper) {
  border-radius: 12px 0 0 12px;
  box-shadow: 0 0 0 1px rgba(0,0,0,0.05) inset;
  padding: 8px 16px;
}
.cdkey-redemption-card :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px var(--theme-color, #FE2C55) inset;
}
.redeem-btn {
  border: none;
  background: transparent;
  color: white;
  font-weight: 500;
  padding: 0 24px;
  height: 100%;
}
.redeem-btn:hover {
  background-color: rgba(0,0,0,0.1);
  color: white;
}

.points-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 32px 16px;
  border-bottom: 1px solid transparent;
}

.points-title {
  font-size: 20px;
  font-weight: 600;
  color: #111827;
}

.points-close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 4px;
  cursor: pointer;
  color: #909399;
  background-color: transparent;
  transition: all 0.3s;
}

.points-close-btn:hover {
  background-color: #f5f7fa;
  color: #303133;
}

.points-balance-card {
  margin: 0 32px 24px;
  padding: 24px;
  border-radius: 12px;
  background: #f3f4f6;
  color: #111827;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: none;
}

.balance-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.balance-label {
  font-size: 14px;
  color: #6b7280;
  font-weight: 500;
}

.balance-value {
  font-size: 32px;
  font-weight: 700;
  color: #111827;
  line-height: 1.1;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.balance-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.balance-recharge-btn,
.balance-redeem-btn {
  height: 28px;
  padding: 0 14px;
  border-radius: 14px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  outline: none;
  white-space: nowrap;
  flex-shrink: 0;
  box-sizing: border-box;
}

.balance-recharge-btn {
  background-color: #111827;
  color: #ffffff;
  border: none;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.balance-recharge-btn:hover {
  background-color: #1f2937;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

.balance-recharge-btn:active {
  transform: translateY(1px);
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.balance-redeem-btn {
  background-color: #ffffff;
  color: #111827;
  border: 1px solid #d1d5db;
}

.balance-redeem-btn:hover {
  background-color: #f9fafb;
  border-color: #9ca3af;
  transform: translateY(-1px);
}

.balance-redeem-btn:active {
  transform: translateY(1px);
  background-color: #f3f4f6;
}



.premium-redeem-container {
  position: relative;
  padding: 40px 32px 32px;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.redeem-close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  cursor: pointer;
  border-radius: 50%;
  transition: all 0.2s;
}
.redeem-close-btn:hover {
  background-color: #f3f4f6;
  color: #111827;
}

.redeem-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}

.redeem-icon-wrapper {
  width: 64px;
  height: 64px;
  background: linear-gradient(135deg, #111827 0%, #374151 100%);
  color: #ffffff;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
  box-shadow: 0 8px 16px rgba(17, 24, 39, 0.15);
  transform: rotate(-3deg);
}

.redeem-title {
  font-size: 22px;
  font-weight: 600;
  color: #111827;
  margin: 0 0 8px;
}

.redeem-desc {
  font-size: 14px;
  color: #6b7280;
  margin: 0 0 28px;
}

.redeem-body {
  width: 100%;
  margin-bottom: 32px;
}

.custom-redeem-input {
  width: 100%;
  height: 52px;
  background: #f9fafb;
  border: 2px solid #e5e7eb;
  border-radius: 12px;
  padding: 0 16px;
  font-size: 16px;
  color: #111827;
  text-align: center;
  letter-spacing: 1px;
  transition: all 0.3s;
  outline: none;
  box-sizing: border-box;
}
.custom-redeem-input:focus {
  border-color: #111827;
  background: #ffffff;
  box-shadow: 0 0 0 4px rgba(17, 24, 39, 0.05);
}
.custom-redeem-input::placeholder {
  color: #9ca3af;
  letter-spacing: normal;
}

.redeem-footer {
  width: 100%;
  display: flex;
  gap: 12px;
}

.redeem-action-btn {
  flex: 1;
  height: 44px;
  border-radius: 10px;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  outline: none;
}

.redeem-action-btn.cancel {
  background: #ffffff;
  color: #374151;
  border: 1px solid #d1d5db;
}
.redeem-action-btn.cancel:hover {
  background: #f9fafb;
  border-color: #9ca3af;
}

.redeem-action-btn.confirm {
  background: #111827;
  color: #ffffff;
  border: none;
  box-shadow: 0 4px 12px rgba(17, 24, 39, 0.15);
}
.redeem-action-btn.confirm:hover:not(:disabled) {
  background: #1f2937;
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(17, 24, 39, 0.2);
}
.redeem-action-btn.confirm:disabled {
  background: #9ca3af;
  cursor: not-allowed;
  box-shadow: none;
}

.points-list-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 0 20px 20px;
  overflow: hidden;
}

.list-section-title {
  font-size: 14px;
  font-weight: 600;
  color: #374151;
  margin-bottom: 12px;
}

.points-list-content {
  flex: 1;
  overflow-y: auto;
  padding-right: 4px;
}

.points-list-content::-webkit-scrollbar {
  width: 4px;
}
.points-list-content::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 2px;
}
.points-list-content::-webkit-scrollbar-track {
  background: transparent;
}

.points-list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f3f4f6;
}

.points-list-item:last-child {
  border-bottom: none;
}

.item-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.item-reason {
  font-size: 14px;
  font-weight: 500;
  color: #111827;
}

.item-reason-detail {
  font-size: 12px;
  color: #6b7280;
  margin-top: -2px;
}

.item-time {
  font-size: 11px;
  color: #9ca3af;
}

.item-right {
  font-size: 15px;
  font-weight: 600;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
}

.item-right.is-positive {
  color: #10b981;
}

.item-right.is-negative {
  color: #111827;
}

.points-pagination-minimal {
  margin-top: 16px;
  display: flex;
  justify-content: center;
}

:deep(.points-pagination-minimal .el-pagination) {
  --el-pagination-bg-color: transparent;
  --el-pagination-hover-color: #111827;
}
.elegant-profile-dropdown {
  padding: 12px !important;
  min-width: 280px;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.6);
  box-shadow: 0 10px 40px -10px rgba(0, 0, 0, 0.1), 0 1px 3px rgba(0, 0, 0, 0.05), inset 0 1px 0 rgba(255, 255, 255, 0.6);
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  box-sizing: border-box;
  z-index: 2001 !important; /* 强制下拉菜单本身也具有超高层级 */
}

.elegant-header {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 8px 16px;
}

.elegant-avatar-wrap {
  position: relative;
  box-shadow: 0 4px 10px rgba(0, 0, 0, 0.08);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.elegant-info {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.elegant-name {
  font-size: 16px;
  font-weight: 600;
  color: #111827;
  line-height: 1.2;
}

.elegant-email {
  font-size: 12px;
  color: #6b7280;
  margin-top: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 140px;
}

.elegant-points-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: linear-gradient(135deg, #fffbeb 0%, #fef3c7 50%, #fce7f3 100%);
  border: 1px solid rgba(251, 191, 36, 0.3);
  box-shadow: 0 4px 12px rgba(251, 191, 36, 0.15), inset 0 1px 0 rgba(255, 255, 255, 0.8);
  border-radius: 12px;
  padding: 12px;
  margin: 0 4px 12px;
  position: relative;
  z-index: 1;
}

.elegant-points-card::before {
  content: '';
  position: absolute;
  top: 0; left: 0; right: 0; bottom: 0;
  background: radial-gradient(circle at top left, rgba(255,255,255,0.8) 0%, transparent 60%);
  border-radius: 12px;
  z-index: 0;
  pointer-events: none;
}

.elegant-points-left, .elegant-points-right {
  position: relative;
  z-index: 2;
}

.elegant-points-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.elegant-points-right {
  display: flex;
  align-items: center;
}

.elegant-points-label {
  font-size: 12px;
  color: #92400e;
  font-weight: 500;
}

.elegant-points-value {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 20px;
  font-weight: 800;
  color: #92400e;
  letter-spacing: -0.5px;
}

.elegant-points-value .el-icon {
  color: #f59e0b;
  filter: drop-shadow(0 2px 4px rgba(245, 158, 11, 0.4));
}

.elegant-recharge-btn {
  background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%) !important;
  color: white !important;
  border: none !important;
  font-weight: 600 !important;
  z-index: 2;
  white-space: nowrap;
  padding: 8px 16px;
  box-shadow: 0 2px 8px rgba(217, 119, 6, 0.3) !important;
  transition: all 0.3s ease !important;
}

.elegant-recharge-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(217, 119, 6, 0.4) !important;
  background: linear-gradient(135deg, #fbbf24 0%, #f59e0b 100%) !important;
}

.elegant-divider {
  height: 1px;
  background: linear-gradient(to right, transparent, rgba(229, 231, 235, 0.8), transparent);
  margin: 8px 0;
}

.elegant-menu-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.elegant-menu-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  color: #4b5563;
  background: transparent !important;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.elegant-icon {
  font-size: 18px;
  margin-right: 12px;
  color: #9ca3af;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.elegant-arrow {
  margin-left: auto;
  font-size: 14px;
  color: #d1d5db;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.elegant-menu-item:hover {
  background: rgba(243, 244, 246, 0.8) !important;
  color: #111827;
  transform: translateX(4px);
}

.elegant-menu-item:hover .elegant-icon {
  color: #3b82f6;
  transform: scale(1.1);
}

.elegant-menu-item:hover .elegant-arrow {
  transform: translateX(2px);
  color: #9ca3af;
}

.elegant-logout:hover {
  background: rgba(254, 242, 242, 0.8) !important;
  color: #ef4444;
}

.elegant-logout:hover .elegant-icon {
  color: #ef4444;
}

/* 确保 teleport 后的下拉菜单在移动端有极高的层级 */
.el-dropdown__popper.el-popper {
  z-index: 3000 !important;
}
.combined-settings-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #f4f4f5;
  border: 1px solid transparent;
  border-radius: 20px;
  padding: 6px 14px;
  min-height: 34px;
  height: auto;
  font-size: 13px;
  color: #3f3f46;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: none;
  flex: 0 0 auto;
  white-space: nowrap;
  overflow: visible; /* 防止文字被上下遮挡 */
}

.combined-settings-btn:hover {
  background: #e4e4e7;
  color: #18181b;
  transform: translateY(-1px);
}

.combined-settings-btn:active {
  transform: translateY(0);
}

.combined-settings-btn .el-icon {
  font-size: 15px;
  color: #71717a;
  transition: color 0.25s ease;
}

.combined-settings-btn:hover .el-icon {
  color: #18181b;
}

.combined-settings-btn .btn-text {
  display: inline-flex;
  align-items: center;
}

.combined-settings-btn .divider {
  width: 1px;
  height: 10px; /* 减小分割线高度 */
  background-color: #e4e4e7; /* 稍微减淡颜色，使其不那么刺眼 */
  margin: 0 6px; /* 减小左右间距 */
}

.combined-settings-btn .unit-text {
  font-size: 12px;
  margin-left: 2px;
  color: #a1a1aa;
  font-weight: normal;
}

/* Settings Popover */
.settings-popover {
  padding: 16px !important;
  border-radius: 12px !important;
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1), 0 8px 10px -6px rgba(0, 0, 0, 0.1) !important;
  border: 1px solid #f3f4f6 !important;
  max-width: 90vw !important;
}

.settings-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.setting-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.setting-label {
  font-size: 13px;
  font-weight: 500;
  color: #374151;
}

.setting-options {
  display: flex;
  flex-wrap: nowrap;
  gap: 8px;
  overflow-x: auto;
  padding-bottom: 4px;
  width: 100%;
  box-sizing: border-box;
}

.setting-options::-webkit-scrollbar {
  display: none;
}

.option-btn {
  flex-shrink: 0;
  padding: 6px 12px;
  border-radius: 6px;
  background: #f3f4f6;
  color: #4b5563;
  font-size: 13px;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid transparent;
}

.option-btn.disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background-color: #f9fafb;
}

.option-btn:not(.disabled):hover {
  background: #e5e7eb;
}

.option-btn.active {
  background: #eff6ff;
  color: #3b82f6;
  border-color: #bfdbfe;
  font-weight: 500;
}

.model-dropdown-menu .el-dropdown-item.is-active-model {
  color: #3b82f6;
  font-weight: 500;
  background-color: #eff6ff;
}

.shimmer-tag {
  position: relative;
  overflow: hidden;
}

.shimmer-tag::after {
  content: "";
  position: absolute;
  top: 0;
  left: -100%;
  width: 50%;
  height: 100%;
  background: linear-gradient(to right, rgba(255, 255, 255, 0) 0%, rgba(255, 255, 255, 0.4) 50%, rgba(255, 255, 255, 0) 100%);
  transform: skewX(-20deg);
  animation: sweep-light 2.5s infinite;
}

@keyframes sweep-light {
  0% {
    left: -100%;
  }
  20% {
    left: 200%;
  }
  100% {
    left: 200%;
  }
}

/* Header Styles */
.main-header {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 64px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 24px;
  background: transparent;
  z-index: 150;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
  flex: 1;
  min-width: 0;
}

.mobile-menu-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  -webkit-app-region: no-drag;
  cursor: pointer;
  padding: 6px;
  border-radius: 8px;
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  color: #4b5563;
  transition: all 0.2s ease;
}

.mobile-menu-btn:hover {
  background-color: #f9fafb;
  border-color: #d1d5db;
  color: #111827;
}

.header-new-chat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background-color: #ffffff;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  color: #4b5563;
  cursor: pointer;
  transition: all 0.2s ease;
}

.header-new-chat-icon:hover {
  background-color: #f9fafb;
  border-color: #d1d5db;
  color: #111827;
}

.header-new-chat-icon .el-icon {
  font-size: 16px;
}

.header-actions {
  flex-shrink: 0;
  margin-left: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-login-btn {
  padding: 6px 16px;
  background-color: #111827;
  color: #ffffff;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.header-login-btn:hover {
  background-color: #374151;
}

.header-points-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  background-color: #ffffff;
  border: 1px solid #e5e7eb;
  color: #374151;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

.header-points-btn:hover {
  background-color: #f9fafb;
  border-color: #d1d5db;
  color: #111827;
}

.header-points-btn .el-icon {
  color: #f59e0b;
  font-size: 15px;
}

.explore-content-area {
  flex: 1;
  overflow-y: auto;
  padding: 84px 32px 32px 32px; /* 顶部留出 header 的空间 */
}

.explore-container {
  max-width: 1400px;
  margin: 0 auto;
}

/* Loading and No More Data styles */
.loading-more, .no-more-data {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30px 0;
  color: #6b7280;
  font-size: 14px;
  gap: 12px;
}

.loading-more .el-icon {
  font-size: 18px;
  color: #3b82f6;
  animation: loading-rotate 2s linear infinite;
}

@keyframes loading-rotate {
  100% {
    transform: rotate(360deg);
  }
}

.no-more-data .divider-line {
  height: 1px;
  width: 40px;
  background-color: #e5e7eb;
}

/* Category Filter Styles */
.category-filter-section {
  width: 100%;
  max-width: 1400px;
  margin: 0 auto 32px;
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 0 16px;
  box-sizing: border-box;
}

.main-categories-wrapper {
  width: 100%;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none; /* Firefox */
  padding-bottom: 4px;
}

.main-categories-wrapper::-webkit-scrollbar {
  display: none; /* Chrome, Safari and Opera */
}

.main-categories {
  display: inline-flex;
  gap: 12px;
  align-items: center;
  padding: 4px;
}

.category-tab {
  padding: 10px 24px;
  border-radius: 100px;
  font-size: 15px;
  font-weight: 600;
  color: #6b7280;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(229, 231, 235, 0.8);
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.02);
  white-space: nowrap;
  user-select: none;
}

.category-tab:hover {
  color: #111827;
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.05);
  background: #ffffff;
}

.category-tab.active {
  background: #111827;
  color: #ffffff;
  border-color: #111827;
  box-shadow: 0 8px 16px rgba(17, 24, 39, 0.15);
  transform: translateY(-2px);
}

.sub-categories {
  width: 100%;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
}

.sub-categories::-webkit-scrollbar {
  display: none;
}

.sub-categories-inner {
  display: inline-flex;
  gap: 10px;
  align-items: center;
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.5);
  border-radius: 20px;
  border: 1px solid rgba(229, 231, 235, 0.6);
  backdrop-filter: blur(12px);
}

.sub-category-tag {
  padding: 6px 18px;
  border-radius: 100px;
  font-size: 14px;
  font-weight: 500;
  color: #6b7280;
  background: transparent;
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
  user-select: none;
}

.sub-category-tag:hover {
  color: #374151;
  background: rgba(243, 244, 246, 0.8);
}

.sub-category-tag.active {
  color: #3b82f6;
  background: #eff6ff;
  font-weight: 600;
}

/* Prompts Gallery 风格瀑布流/网格 */
.explore-masonry {
  columns: 4;
  column-gap: 20px;
  width: 100%;
}

@media (max-width: 1200px) {
  .explore-masonry {
    columns: 3;
  }
}

@media (max-width: 900px) {
  .explore-masonry {
    columns: 2;
  }
}

@media (max-width: 600px) {
  .explore-masonry {
    columns: 1;
  }
}

/* 提示词卡片样式 (类似 lexica.art / promptsref) */
.prompt-card {
  break-inside: avoid;
  margin-bottom: 20px;
  position: relative;
  border-radius: 16px;
  overflow: hidden;
  cursor: pointer;
  background-color: #f3f4f6;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: block; /* 防止 inline-block 导致的高度塌陷 */
}

.prompt-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
}

.card-image {
  width: 100%;
  display: block; /* 移除底部空白 */
  height: auto; /* 允许高度自适应，形成真实的瀑布流 */
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

.prompt-card:hover .card-image {
  transform: scale(1.08);
}

.explore-no-image, .image-placeholder, .image-error {
  width: 100%;
  aspect-ratio: 3 / 4; /* 为没有加载出图片的状态保留一个基础高度 */
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  font-size: 32px;
}

/* 悬浮遮罩和内容 */
.card-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(to top, rgba(17, 24, 39, 0.9) 0%, rgba(17, 24, 39, 0.4) 40%, transparent 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 20px;
}

.prompt-card:hover .card-overlay {
  opacity: 1;
}

.overlay-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.card-prompt-text {
  color: #ffffff;
  font-size: 14px;
  line-height: 1.5;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
  text-shadow: 0 2px 4px rgba(0,0,0,0.5);
  transform: translateY(10px);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.prompt-card:hover .card-prompt-text {
  transform: translateY(0);
}

.overlay-actions {
  display: flex;
  justify-content: flex-end;
  transform: translateY(10px);
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transition-delay: 0.05s;
}

.prompt-card:hover .overlay-actions {
  transform: translateY(0);
}

.try-btn {
  background: rgba(255, 255, 255, 0.2) !important;
  color: #ffffff !important;
  border: 1px solid rgba(255, 255, 255, 0.4) !important;
  backdrop-filter: blur(8px);
  font-weight: 500;
  padding: 8px 20px;
  transition: all 0.2s ease;
}

.try-btn:hover {
  background: #ffffff !important;
  color: #111827 !important;
  transform: scale(1.05);
}

.explore-custom-input :deep(.el-textarea__inner) {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
  padding: 0 !important;
  font-size: 16px;
  line-height: 1.6;
  color: #374151;
}

.explore-custom-input :deep(.el-textarea__inner:focus) {
  box-shadow: none !important;
}

@media (min-width: 769px) {
  .mobile-menu-btn {
    display: none;
  }
  .mobile-only {
    display: none !important;
  }
}

@media (max-width: 768px) {
  .desktop-only {
    display: none !important;
  }
  .main-header {
    padding: 0 16px;
  }
  .explore-content-area {
    padding: 84px 16px 16px 16px;
    overflow-x: hidden;
  }
  .explore-masonry {
    columns: 2; /* 移动端下保持两列 */
    column-gap: 12px;
  }
  .prompt-card {
    margin-bottom: 12px;
  }
  .mobile-menu-btn {
    font-size: 20px;
    -webkit-tap-highlight-color: transparent;
    outline: none;
  }
}
</style>
