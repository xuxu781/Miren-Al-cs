<template>
  <div class="admin-container fade-in">
    <div class="page-header">
      <div class="header-info">
        <div class="icon-wrapper">
          <el-icon :size="24" color="#409eff"><Cpu /></el-icon>
        </div>
        <div class="header-title">
          <h2>模型上游</h2>
          <span class="subtitle">统一配置模型连接、能力、部署、计费和生成执行策略。支持一键切换主用上游。</span>
        </div>
      </div>
      <div class="header-actions">
        <el-input
          v-model="searchQuery"
          placeholder="搜索上游模型 ID 或提供商..."
          :prefix-icon="Search"
          clearable
          class="search-input"
        />
        <el-button type="primary" class="add-btn" @click="showAddDialog = true">
          <el-icon class="el-icon--left"><Plus /></el-icon>
          添加上游
        </el-button>
      </div>
    </div>

    <el-card shadow="never" class="admin-card">
      <div class="table-toolbar">
        <el-radio-group v-model="filterType" class="type-filter-group">
          <el-radio-button :value="null">全部</el-radio-button>
          <el-radio-button :value="1">文本生成</el-radio-button>
          <el-radio-button :value="2">视频生成</el-radio-button>
          <el-radio-button :value="4">图片生成</el-radio-button>
        </el-radio-group>
      </div>

      <el-table 
        :data="pagedList" 
        v-loading="loading" 
        style="width: 100%" 
        :header-cell-style="{ background: '#f8fafc', color: '#475569', fontWeight: '600', height: '54px' }"
        row-class-name="custom-table-row"
      >
        <el-table-column prop="model_type" label="模型类型" width="100" align="center">
          <template #default="scope">
            <el-tag :type="getModelTypeTag(scope.row.model_type)" size="small">
              {{ getModelTypeName(scope.row.model_type) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="series_id" label="系列 ID(显示名称)" min-width="160">
          <template #default="scope">
            <div class="font-bold text-gray-800">{{ scope.row.series_id || '-' }}</div>
          </template>
        </el-table-column>

        <el-table-column prop="logical_model" label="上游模型 ID" min-width="120">
          <template #default="scope">
            <div class="font-bold text-gray-800">{{ scope.row.logical_model }}</div>
          </template>
        </el-table-column>
        
        <el-table-column prop="provider" label="提供商" min-width="120">
          <template #default="scope">
            <el-tag size="small" type="info">{{ scope.row.provider }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="channel_name" label="渠道名称" min-width="120">
          <template #default="scope">
            <el-tag size="small" type="warning" v-if="scope.row.channel_name">{{ scope.row.channel_name }}</el-tag>
            <span v-else class="text-gray-400 text-xs">-</span>
          </template>
        </el-table-column>

        <el-table-column prop="connection_url" label="连接地址" min-width="180">
          <template #default="scope">
            <div class="truncate text-sm text-gray-500" :title="scope.row.connection_url">
              {{ scope.row.connection_url }}
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="is_primary" label="主用上游" width="100" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.is_primary ? 'success' : 'info'" size="small" effect="dark">
              {{ scope.row.is_primary ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'danger'" size="small" effect="light" class="status-tag">
              <span class="status-dot" :class="scope.row.status === 'active' ? 'dot-success' : 'dot-danger'"></span>
              {{ scope.row.status === 'active' ? '正常' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="220" align="center" fixed="right">
          <template #default="scope">
            <div class="action-btns">
              <el-tooltip content="设为主用" placement="top" v-if="!scope.row.is_primary">
                <el-button type="success" link @click="handleSetPrimary(scope.row)" class="action-btn">
                  <el-icon :size="16"><Check /></el-icon>
                </el-button>
              </el-tooltip>

              <el-tooltip content="编辑上游配置" placement="top">
                <el-button type="primary" link @click="openEditDialog(scope.row)" class="action-btn">
                  <el-icon :size="16"><EditPen /></el-icon>
                </el-button>
              </el-tooltip>
              
              <el-tooltip content="删除上游" placement="top">
                <div style="display: inline-block;">
                  <el-popconfirm 
                    title="确定要删除该上游配置吗？此操作不可恢复！" 
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
          <el-empty description="暂无模型上游数据" :image-size="120" />
        </template>
      </el-table>

      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :background="true"
          layout="total, sizes, prev, pager, next, jumper"
          :total="filteredList.length"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 添加上游对话框 -->
    <el-dialog v-model="showAddDialog" title="添加模型上游" width="600px" destroy-on-close class="custom-dialog">
      <div class="dialog-desc">填写模型上游的连接和配置信息。</div>
      <el-form ref="addFormRef" :model="addForm" :rules="rules" label-width="auto" size="large" label-position="right">
        <el-form-item label="模型类型" prop="model_type">
          <el-select v-model="addForm.model_type" placeholder="请选择模型类型" style="width: 100%" @change="handleModelTypeChange('add')">
            <el-option label="文本生成" :value="1" />
            <el-option label="视频生成" :value="2" />
            <el-option label="图片生成" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="系列 ID(显示名称)" prop="series_id">
          <el-input v-model="addForm.series_id" placeholder="供用户调用的模型系列也是前端展示名称，例如: gpt-4" />
        </el-form-item>
        <el-form-item label="上游模型 ID" prop="logical_model">
          <el-input v-model="addForm.logical_model" placeholder="后端实际调用的模型，例如: dall-e-3" />
        </el-form-item>
        <el-form-item label="提供商" prop="provider">
          <el-select v-model="addForm.provider" placeholder="请选择提供商" filterable allow-create style="width: 100%">
            <el-option v-for="item in getProviderOptions(addForm.model_type)" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="渠道名称" prop="channel_name">
          <el-input v-model="addForm.channel_name" placeholder="标识当前第三方API来源渠道，例如: 中转A站、直连等" />
        </el-form-item>
        <el-form-item label="连接地址" prop="connection_url">
          <el-input v-model="addForm.connection_url" placeholder="API Base URL" />
        </el-form-item>
        <el-form-item label="API Key" prop="api_key">
          <el-input v-model="addForm.api_key" type="password" placeholder="API 密钥" show-password />
        </el-form-item>
        
        <template v-if="addForm.model_type === 4">
          <div @click="showAddImageConfig = !showAddImageConfig" style="cursor: pointer; user-select: none;">
            <el-divider content-position="left">
              <span style="display: flex; align-items: center; gap: 4px; color: #409eff;">
                图片模型配置 (可选)
                <el-icon><ArrowDown v-if="showAddImageConfig" /><ArrowRight v-else /></el-icon>
              </span>
            </el-divider>
          </div>
          
          <el-collapse-transition>
            <div v-show="showAddImageConfig">
              <div v-for="(tier, index) in addForm.resolution_tiers" :key="'add_' + index" style="padding-left: 20px; background: #fafafa; padding-top: 15px; padding-bottom: 5px; border-radius: 8px; margin-bottom: 15px; position: relative;">
            <el-button type="danger" link style="position: absolute; right: 10px; top: 10px; z-index: 1;" @click="handleRemoveResolution(index, 'add')">
              <el-icon><Delete /></el-icon>
            </el-button>

            <div @click="addExpandedTiers[index] = !addExpandedTiers[index]" style="cursor: pointer; color: #606266; font-size: 14px; font-weight: 500; margin-bottom: 15px; display: inline-flex; align-items: center; gap: 4px; user-select: none;">
              <span>档位配置 {{ tier ? `(${tier})` : '' }}</span>
              <el-icon><ArrowDown v-if="addExpandedTiers[index]" /><ArrowRight v-else /></el-icon>
            </div>

            <el-collapse-transition>
              <div v-show="addExpandedTiers[index]">
                <el-form-item label="分辨率档位" :prop="'resolution_tiers.' + index" :rules="[{ required: true, message: '请选择分辨率', trigger: 'change' }]">
              <el-select
                v-model="addForm.resolution_tiers[index]"
                filterable
                allow-create
                default-first-option
                placeholder="请选择或输入分辨率，如1K、2K、4K"
                style="width: calc(100% - 30px)"
                @change="(val: string) => handleResolutionChange(val, index, 'add')"
              >
                <el-option label="1K" value="1K" :disabled="addForm.resolution_tiers.includes('1K') && addForm.resolution_tiers[index] !== '1K'" />
                <el-option label="2K" value="2K" :disabled="addForm.resolution_tiers.includes('2K') && addForm.resolution_tiers[index] !== '2K'" />
                <el-option label="4K" value="4K" :disabled="addForm.resolution_tiers.includes('4K') && addForm.resolution_tiers[index] !== '4K'" />
                <el-option label="8K" value="8K" :disabled="addForm.resolution_tiers.includes('8K') && addForm.resolution_tiers[index] !== '8K'" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="对应模型" :prop="'resolution_configs.' + tier + '.model'">
              <el-input v-model="addForm.resolution_configs[tier].model" placeholder="后端实际调用的模型，为空则使用上方主模型" />
            </el-form-item>
            <el-form-item label="状态" :prop="'resolution_configs.' + tier + '.enabled'">
              <el-select v-model="addForm.resolution_configs[tier].enabled" placeholder="请选择状态" style="width: 120px;">
                <el-option label="启用" :value="true" />
                <el-option label="禁用" :value="false" />
                <el-option label="维护中" value="maintenance" />
              </el-select>
            </el-form-item>
            <el-form-item label="活动标签" :prop="'resolution_configs.' + tier + '.tag'">
              <div style="display: flex; gap: 10px; width: 100%; align-items: center;">
                <el-input v-model="addForm.resolution_configs[tier].tag" placeholder="如：限时免费、5折优惠 (展示在该分辨率选项中)" style="flex: 1;" />
                <span style="font-size: 13px; color: #606266; white-space: nowrap;">背景色:</span>
                <el-color-picker v-model="addForm.resolution_configs[tier].tag_color" :predefine="predefineColors" />
              </div>
            </el-form-item>
            <el-form-item label="积分/张" :prop="'resolution_configs.' + tier + '.credits_per_image'">
              <el-input-number v-model="addForm.resolution_configs[tier].credits_per_image" :min="0" :precision="2" :step="0.1" controls-position="right" style="width: 100%" />
            </el-form-item>
            <el-form-item label="尺寸参数" :prop="'resolution_configs.' + tier + '.size_parameter'">
              <el-select v-model="addForm.resolution_configs[tier].size_parameter" placeholder="请选择尺寸参数类型">
                <el-option label="宽高比(16:9)" value="aspect_ratio" />
                <el-option label="按分辨率转像素" value="pixel_resolution" />
                <el-option label="支持像素尺寸" value="pixel_size" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="支持比例" :prop="'resolution_configs.' + tier + '.aspect_ratios'" :rules="[{ validator: validateAspectRatios, trigger: 'blur' }]">
              <el-input v-model="addForm.resolution_configs[tier].aspect_ratios" placeholder="以逗号分隔，如 auto,1:1,4:3,3:4,3:2,2:3,16:9,9:16,21:9,9:21" />
            </el-form-item>

            <el-form-item label="支持张数" :prop="'resolution_configs.' + tier + '.image_counts'">
              <el-input-number v-model="addForm.resolution_configs[tier].image_counts" :min="1" />
            </el-form-item>
            
            <el-form-item label="最多参考图" :prop="'resolution_configs.' + tier + '.max_reference_images'">
              <el-input-number v-model="addForm.resolution_configs[tier].max_reference_images" :min="0" />
            </el-form-item>
              </div>
            </el-collapse-transition>
          </div>
          
          <div class="add-config-btn-wrapper" v-if="addForm.model_type === 4">
            <el-button type="primary" plain class="add-config-btn" @click="handleAddResolution('add')">
              <el-icon><Plus /></el-icon>
              <span>增加模型配置</span>
            </el-button>
          </div>
            </div>
          </el-collapse-transition>
        </template>
        
        <el-form-item label="模型活动标签" prop="activity_tag" v-if="addForm.model_type !== 1">
          <div style="display: flex; gap: 10px; width: 100%; align-items: center;">
            <el-input v-model="addForm.activity_tag" placeholder="如：限时免费、5折优惠 (展示在模型下拉列表中)" style="flex: 1;" />
            <span style="font-size: 13px; color: #606266; white-space: nowrap;">背景色:</span>
            <el-color-picker v-model="addForm.activity_tag_color" :predefine="predefineColors" />
          </div>
        </el-form-item>

        <el-divider content-position="left">高级配置</el-divider>
        <el-form-item label="超时时间(秒)" prop="timeout_seconds">
          <el-input-number v-model="addForm.timeout_seconds" :min="1" :step="10" placeholder="如 60" />
        </el-form-item>
        <el-form-item label="计费策略" prop="billing_strategy">
          <el-input v-model="addForm.billing_strategy" placeholder="JSON格式计费策略 (可选)" />
        </el-form-item>
        <el-form-item label="主用上游" prop="is_primary">
          <el-switch v-model="addForm.is_primary" active-text="是" inactive-text="否" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="addForm.status" active-value="active" inactive-value="inactive" active-text="正常" inactive-text="禁用" inline-prompt />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddDialog = false">取 消</el-button>
          <el-button type="primary" @click="handleAdd" :loading="submitLoading">确认添加</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑上游对话框 -->
    <el-dialog v-model="showEditDialog" title="编辑模型上游" width="600px" destroy-on-close class="custom-dialog">
      <div class="dialog-desc">修改模型上游配置信息。</div>
      <el-form ref="editFormRef" :model="editForm" :rules="rules" label-width="auto" size="large" label-position="right">
        <el-form-item label="模型类型" prop="model_type">
          <el-select v-model="editForm.model_type" placeholder="请选择模型类型" style="width: 100%" @change="handleModelTypeChange('edit')">
            <el-option label="文本生成" :value="1" />
            <el-option label="视频生成" :value="2" />
            <el-option label="图片生成" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="系列 ID(显示名称)" prop="series_id">
          <el-input v-model="editForm.series_id" placeholder="供用户调用的模型系列也是前端展示名称，例如: gpt-4" />
        </el-form-item>
        <el-form-item label="上游模型 ID" prop="logical_model">
          <el-input v-model="editForm.logical_model" placeholder="后端实际调用的模型，例如: dall-e-3" />
        </el-form-item>
        <el-form-item label="提供商" prop="provider">
          <el-select v-model="editForm.provider" placeholder="请选择提供商" filterable allow-create style="width: 100%">
            <el-option v-for="item in getProviderOptions(editForm.model_type)" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="渠道名称" prop="channel_name">
          <el-input v-model="editForm.channel_name" placeholder="标识当前第三方API来源渠道，例如: 中转A站、直连等" />
        </el-form-item>
        <el-form-item label="连接地址" prop="connection_url">
          <el-input v-model="editForm.connection_url" placeholder="API Base URL" />
        </el-form-item>
        <el-form-item label="API Key" prop="api_key">
          <el-input v-model="editForm.api_key" type="password" placeholder="API 密钥 (留空不修改)" show-password />
        </el-form-item>

        <template v-if="editForm.model_type === 4">
          <div @click="showEditImageConfig = !showEditImageConfig" style="cursor: pointer; user-select: none;">
            <el-divider content-position="left">
              <span style="display: flex; align-items: center; gap: 4px; color: #409eff;">
                图片模型配置 (可选)
                <el-icon><ArrowDown v-if="showEditImageConfig" /><ArrowRight v-else /></el-icon>
              </span>
            </el-divider>
          </div>
          
          <el-collapse-transition>
            <div v-show="showEditImageConfig">
              <div v-for="(tier, index) in editForm.resolution_tiers" :key="'edit_' + index" style="padding-left: 20px; background: #fafafa; padding-top: 15px; padding-bottom: 5px; border-radius: 8px; margin-bottom: 15px; position: relative;">
            <el-button type="danger" link style="position: absolute; right: 10px; top: 10px; z-index: 1;" @click="handleRemoveResolution(index, 'edit')">
              <el-icon><Delete /></el-icon>
            </el-button>

            <div @click="editExpandedTiers[index] = !editExpandedTiers[index]" style="cursor: pointer; color: #606266; font-size: 14px; font-weight: 500; margin-bottom: 15px; display: inline-flex; align-items: center; gap: 4px; user-select: none;">
              <span>档位配置 {{ tier ? `(${tier})` : '' }}</span>
              <el-icon><ArrowDown v-if="editExpandedTiers[index]" /><ArrowRight v-else /></el-icon>
            </div>

            <el-collapse-transition>
              <div v-show="editExpandedTiers[index]">
                <el-form-item label="分辨率档位" :prop="'resolution_tiers.' + index" :rules="[{ required: true, message: '请选择分辨率', trigger: 'change' }]">
              <el-select
                v-model="editForm.resolution_tiers[index]"
                filterable
                allow-create
                default-first-option
                placeholder="请选择或输入分辨率，如1K、2K、4K"
                style="width: calc(100% - 30px)"
                @change="(val: string) => handleResolutionChange(val, index, 'edit')"
              >
                <el-option label="1K" value="1K" :disabled="editForm.resolution_tiers.includes('1K') && editForm.resolution_tiers[index] !== '1K'" />
                <el-option label="2K" value="2K" :disabled="editForm.resolution_tiers.includes('2K') && editForm.resolution_tiers[index] !== '2K'" />
                <el-option label="4K" value="4K" :disabled="editForm.resolution_tiers.includes('4K') && editForm.resolution_tiers[index] !== '4K'" />
                <el-option label="8K" value="8K" :disabled="editForm.resolution_tiers.includes('8K') && editForm.resolution_tiers[index] !== '8K'" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="对应模型" :prop="'resolution_configs.' + tier + '.model'">
              <el-input v-model="editForm.resolution_configs[tier].model" placeholder="后端实际调用的模型，为空则使用上方主模型" />
            </el-form-item>
            <el-form-item label="状态" :prop="'resolution_configs.' + tier + '.enabled'">
              <el-select v-model="editForm.resolution_configs[tier].enabled" placeholder="请选择状态" style="width: 120px;">
                <el-option label="启用" :value="true" />
                <el-option label="禁用" :value="false" />
                <el-option label="维护中" value="maintenance" />
              </el-select>
            </el-form-item>
            <el-form-item label="活动标签" :prop="'resolution_configs.' + tier + '.tag'">
              <div style="display: flex; gap: 10px; width: 100%; align-items: center;">
                <el-input v-model="editForm.resolution_configs[tier].tag" placeholder="如：限时免费、5折优惠 (展示在该分辨率选项中)" style="flex: 1;" />
                <span style="font-size: 13px; color: #606266; white-space: nowrap;">背景色:</span>
                <el-color-picker v-model="editForm.resolution_configs[tier].tag_color" :predefine="predefineColors" />
              </div>
            </el-form-item>
            <el-form-item label="积分/张" :prop="'resolution_configs.' + tier + '.credits_per_image'">
              <el-input-number v-model="editForm.resolution_configs[tier].credits_per_image" :min="0" :precision="2" :step="0.1" controls-position="right" style="width: 100%" />
            </el-form-item>
            <el-form-item label="尺寸参数" :prop="'resolution_configs.' + tier + '.size_parameter'">
              <el-select v-model="editForm.resolution_configs[tier].size_parameter" placeholder="请选择尺寸参数类型">
                <el-option label="宽高比(16:9)" value="aspect_ratio" />
                <el-option label="按分辨率转像素" value="pixel_resolution" />
                <el-option label="支持像素尺寸" value="pixel_size" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="支持比例" :prop="'resolution_configs.' + tier + '.aspect_ratios'" :rules="[{ validator: validateAspectRatios, trigger: 'blur' }]">
              <el-input v-model="editForm.resolution_configs[tier].aspect_ratios" placeholder="以逗号分隔，如 auto,1:1,4:3,3:4,3:2,2:3,16:9,9:16,21:9,9:21" />
            </el-form-item>

            <el-form-item label="支持张数" :prop="'resolution_configs.' + tier + '.image_counts'">
              <el-input-number v-model="editForm.resolution_configs[tier].image_counts" :min="1" />
            </el-form-item>
            
            <el-form-item label="最多参考图" :prop="'resolution_configs.' + tier + '.max_reference_images'">
              <el-input-number v-model="editForm.resolution_configs[tier].max_reference_images" :min="0" />
            </el-form-item>
              </div>
            </el-collapse-transition>
          </div>
          
          <div class="add-config-btn-wrapper" v-if="editForm.model_type === 4">
            <el-button type="primary" plain class="add-config-btn" @click="handleAddResolution('edit')">
              <el-icon><Plus /></el-icon>
              <span>增加模型配置</span>
            </el-button>
          </div>
            </div>
          </el-collapse-transition>
        </template>
        
        <el-form-item label="模型活动标签" prop="activity_tag" v-if="editForm.model_type !== 1">
          <div style="display: flex; gap: 10px; width: 100%; align-items: center;">
            <el-input v-model="editForm.activity_tag" placeholder="如：限时免费、5折优惠 (展示在模型下拉列表中)" style="flex: 1;" />
            <span style="font-size: 13px; color: #606266; white-space: nowrap;">背景色:</span>
            <el-color-picker v-model="editForm.activity_tag_color" :predefine="predefineColors" />
          </div>
        </el-form-item>

        <el-divider content-position="left">高级配置</el-divider>
        <el-form-item label="超时时间(秒)" prop="timeout_seconds">
          <el-input-number v-model="editForm.timeout_seconds" :min="1" :step="10" placeholder="如 60" />
        </el-form-item>
        <el-form-item label="计费策略" prop="billing_strategy">
          <el-input v-model="editForm.billing_strategy" placeholder="JSON格式计费策略 (可选)" />
        </el-form-item>
        <el-form-item label="主用上游" prop="is_primary">
          <el-switch v-model="editForm.is_primary" active-text="是" inactive-text="否" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="editForm.status" active-value="active" inactive-value="inactive" active-text="正常" inactive-text="禁用" inline-prompt />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取 消</el-button>
          <el-button type="primary" @click="handleEdit" :loading="editLoading">确认保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, EditPen, Delete, Search, Warning, Cpu, Check, ArrowDown, ArrowRight } from '@element-plus/icons-vue'

const showAddImageConfig = ref(true)
const showEditImageConfig = ref(true)
const addExpandedTiers = ref<Record<number, boolean>>({})
const editExpandedTiers = ref<Record<number, boolean>>({})

const predefineColors = ref([
  '#10b981', // 绿色
  '#f59e0b', // 橙色
  '#ef4444', // 红色
  '#3b82f6', // 蓝色
  '#8b5cf6', // 紫色
  '#ec4899', // 粉色
  '#14b8a6', // 翠绿
  '#f97316', // 亮橙
])

const upstreamList = ref<any[]>([])
const loading = ref(false)

// 搜索和分页
const searchQuery = ref('')
const filterType = ref<number | null>(null)
const currentPage = ref(1)
const pageSize = ref(10)

const filteredList = computed(() => {
  let list = upstreamList.value
  
  if (filterType.value !== null && filterType.value !== undefined && String(filterType.value) !== '') {
    list = list.filter(item => item.model_type === Number(filterType.value))
  }
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    list = list.filter(item => 
      (item.logical_model && item.logical_model.toLowerCase().includes(query)) ||
      (item.provider && item.provider.toLowerCase().includes(query))
    )
  }
  
  return list
})

const pagedList = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredList.value.slice(start, end)
})

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
}

const showAddDialog = ref(false)
const submitLoading = ref(false)
const addFormRef = ref()
const addForm = ref({
  model_type: 4,
  series_id: '',
  logical_model: '',
  provider: '',
  channel_name: '',
  connection_url: '',
  api_key: '',
  display_name: '',
  resolution_tiers: [] as string[],
  resolution_configs: {} as Record<string, any>,
  aspect_ratios: 'auto,1:1,4:3,3:4,3:2,2:3,16:9,9:16,21:9,9:21',
  image_counts: 4,
  activity_tag: '',
  activity_tag_color: '#10b981',
  size_parameter: 'aspect_ratio',
  max_reference_images: 3,
  timeout_seconds: 60,
  billing_strategy: '',
  is_primary: false,
  status: 'active'
})

// Deprecated

const handleAddResolution = (type: 'add' | 'edit') => {
  const form = type === 'add' ? addForm.value : editForm.value
  
  if (form.resolution_tiers.some(t => !t || t.trim() === '')) {
    ElMessage.warning('请先选择或输入当前的分辨率档位名称')
    return
  }
  
  const newTier = ''
  
  form.resolution_tiers.push(newTier)
  form.resolution_configs[newTier] = {
    model: '',
    tag: '',
    tag_color: '#10b981',
    enabled: true,
    aspect_ratios: 'auto,1:1,4:3,3:4,3:2,2:3,16:9,9:16,21:9,9:21',
    image_counts: 4,
    credits_per_image: 1.0,
    size_parameter: 'aspect_ratio',
    max_reference_images: 3
  }
  
  const newIndex = form.resolution_tiers.length - 1
  if (type === 'add') {
    addExpandedTiers.value[newIndex] = true
  } else {
    editExpandedTiers.value[newIndex] = true
  }
}

const handleRemoveResolution = (index: number, type: 'add' | 'edit') => {
  const form = type === 'add' ? addForm.value : editForm.value
  const tier = form.resolution_tiers[index]
  form.resolution_tiers.splice(index, 1)
  delete form.resolution_configs[tier]
}

const handleResolutionChange = (newVal: string, index: number, type: 'add' | 'edit') => {
  const form = type === 'add' ? addForm.value : editForm.value
  
  // 检查是否选择了重复的值（忽略大小写，防止输入 1k 绕过 1K 的限制）
  const upperVal = newVal.toUpperCase()
  const duplicateCount = form.resolution_tiers.filter((t: string) => t.toUpperCase() === upperVal).length
  
  if (duplicateCount > 1) {
    ElMessage.warning('该分辨率档位已存在，请选择或输入其他名称')
    // 找出原来的值（存在于 configs 中但不在当前 tiers 中的 key）
    const oldTierKey = Object.keys(form.resolution_configs).find(t => !form.resolution_tiers.includes(t)) || ''
    // 延迟恢复原来的值以覆盖组件内部状态
    setTimeout(() => {
      form.resolution_tiers.splice(index, 1, oldTierKey)
    }, 50)
    return
  }

  // 自动将小写字母转换为大写字母
  if (newVal !== upperVal) {
    setTimeout(() => {
      form.resolution_tiers.splice(index, 1, upperVal)
      if (form.resolution_configs[newVal]) {
        form.resolution_configs[upperVal] = { ...form.resolution_configs[newVal] }
        delete form.resolution_configs[newVal]
      }
    }, 10)
  }

  const oldTiers = Object.keys(form.resolution_configs)
  const missingTier = form.resolution_tiers.find(t => !oldTiers.includes(t)) || newVal
  const oldTierKey = oldTiers.find(t => !form.resolution_tiers.includes(t))
  
  if (oldTierKey !== undefined && oldTierKey !== missingTier) {
    form.resolution_configs[missingTier] = { ...form.resolution_configs[oldTierKey] }
    delete form.resolution_configs[oldTierKey]
  }
}

const showEditDialog = ref(false)
const editLoading = ref(false)
const editFormRef = ref()
const editForm = ref({
  id: null as number | null,
  model_type: 4,
  series_id: '',
  logical_model: '',
  provider: '',
  channel_name: '',
  connection_url: '',
  api_key: '',
  display_name: '',
  resolution_tiers: [] as string[],
  resolution_configs: {} as Record<string, any>,
  aspect_ratios: 'auto,1:1,4:3,3:4,3:2,2:3,16:9,9:16,21:9,9:21',
  image_counts: 4,
  activity_tag: '',
  activity_tag_color: '#10b981',
  size_parameter: 'aspect_ratio',
  max_reference_images: 3,
  timeout_seconds: 60,
  billing_strategy: '',
  is_primary: false,
  status: 'active'
})

// Deprecated

// Removed unused computed availableUpstreams

const validateAspectRatios = (_rule: any, value: string, callback: any) => {
  if (!value) {
    return callback()
  }
  const regex = /^(auto|\d+:\d+)(,(auto|\d+:\d+))*$/
  if (!regex.test(value)) {
    callback(new Error('比例格式不正确，请使用英文逗号分隔，如 auto,1:1,16:9'))
  } else {
    callback()
  }
}

const rules = {
  model_type: [{ required: true, message: '请选择模型类型', trigger: 'change' }],
  logical_model: [{ required: true, message: '请输入上游模型 ID', trigger: 'blur' }],
  provider: [{ required: true, message: '请输入提供商', trigger: 'blur' }],
  channel_name: [{ required: true, message: '请输入渠道名称', trigger: 'blur' }],
  connection_url: [{ required: true, message: '请输入连接地址', trigger: 'blur' }],
  aspect_ratios: [{ validator: validateAspectRatios, trigger: 'blur' }]
}

const getModelTypeName = (type: number) => {
  switch (type) {
    case 1: return '文本生成'
    case 2: return '视频生成'
    case 4: return '图片生成'
    default: return '未知'
  }
}

const getModelTypeTag = (type: number) => {
  switch (type) {
    case 1: return ''
    case 2: return 'warning'
    case 4: return 'success'
    default: return 'info'
  }
}

const getProviderOptions = (modelType: number) => {
  switch (modelType) {
    case 1: // 文本生成
      return [
        { label: 'OpenAI (文本标准协议)', value: 'openai_text' },
        { label: 'Gemini (文本原生协议)', value: 'gemini_text' },
        { label: 'Claude (Anthropic协议)', value: 'claude_text' }
      ]
    case 2: // 视频生成
      return [
        { label: 'OpenAI (视频标准协议)', value: 'openai_video' },
        { label: 'Runway (视频协议)', value: 'runway_video' },
        { label: 'Kling (可灵视频协议)', value: 'kling_video' }
      ]
    case 4: // 图片生成
      return [
        { label: 'OpenAI (图片标准协议)', value: 'openai_image' },
        { label: 'Gemini (图片原生协议)', value: 'gemini_image' }
      ]
    default:
      return [
        { label: 'OpenAI (标准协议)', value: 'openai' },
        { label: 'Gemini (Google原生协议)', value: 'gemini' }
      ]
  }
}

const handleModelTypeChange = (type: 'add' | 'edit') => {
  const form = type === 'add' ? addForm.value : editForm.value
  const options = getProviderOptions(form.model_type)
  const isExist = options.some(opt => opt.value === form.provider)
  if (!isExist && options.length > 0) {
    form.provider = options[0].value
  } else if (!isExist) {
    form.provider = ''
  }
}

const getHeaders = () => {
  const token = localStorage.getItem('token')
  return {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  }
}

const fetchUpstreamList = async () => {
  loading.value = true
  try {
    const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    const res = await fetch(`${apiUrl}/api/admin/upstreams/list`, { headers: getHeaders() })
    const data = await res.json()
    if (res.ok && data.success) {
      upstreamList.value = data.data || []
    } else {
      ElMessage.error(data.error || '获取上游列表失败')
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
        const payload: any = { ...addForm.value }
        if (Array.isArray(payload.resolution_tiers)) {
          payload.resolution_tiers = payload.resolution_tiers.join(',')
        }
        if (payload.image_counts !== undefined) {
          payload.image_counts = String(payload.image_counts)
        }
        if (payload.resolution_configs) {
          payload.resolution_configs = JSON.stringify(payload.resolution_configs)
        }
        
        const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
        const res = await fetch(`${apiUrl}/api/admin/upstreams/create`, {
          method: 'POST',
          headers: getHeaders(),
          body: JSON.stringify(payload)
        })
        const data = await res.json()
        if (res.ok && data.success) {
          ElMessage.success('添加成功')
          showAddDialog.value = false
          addForm.value = { model_type: 4, series_id: '', logical_model: '', provider: '', channel_name: '', connection_url: '', api_key: '', display_name: '', resolution_tiers: [], resolution_configs: {}, aspect_ratios: 'auto,1:1,4:3,3:4,3:2,2:3,16:9,9:16,21:9,9:21', image_counts: 4, activity_tag: '', activity_tag_color: '#10b981', size_parameter: 'aspect_ratio', max_reference_images: 3, timeout_seconds: 60, billing_strategy: '', is_primary: false, status: 'active' }
          fetchUpstreamList()
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
  let parsedConfigs: Record<string, any> = {}
  if (row.resolution_configs) {
    try {
      parsedConfigs = JSON.parse(row.resolution_configs)
      for (const key in parsedConfigs) {
        if (parsedConfigs[key].enabled === undefined) {
          parsedConfigs[key].enabled = true
        }
      }
    } catch (e) {
      console.error('Failed to parse resolution_configs', e)
    }
  }

  editForm.value = {
    id: row.id,
    model_type: row.model_type,
    series_id: row.series_id,
    logical_model: row.logical_model,
    provider: row.provider,
    channel_name: row.channel_name || '',
    connection_url: row.connection_url,
    api_key: row.api_key,
    display_name: row.display_name || '',
    resolution_tiers: row.resolution_tiers ? row.resolution_tiers.split(',').filter((t: string) => t) : [],
    resolution_configs: parsedConfigs,
    aspect_ratios: row.aspect_ratios || 'auto,1:1,4:3,3:4,3:2,2:3,16:9,9:16,21:9,9:21',
    image_counts: row.image_counts && !isNaN(Number(row.image_counts)) ? Number(row.image_counts) : 4,
    activity_tag: row.activity_tag || '',
    activity_tag_color: row.activity_tag_color || '#10b981',
    size_parameter: row.size_parameter || 'aspect_ratio',
    max_reference_images: row.max_reference_images || 0,
    timeout_seconds: row.timeout_seconds || 60,
    billing_strategy: row.billing_strategy,
    is_primary: row.is_primary,
    status: row.status
  }
  showEditDialog.value = true
}

const handleEdit = async () => {
  if (!editFormRef.value) return
  await editFormRef.value.validate(async (valid: boolean) => {
    if (valid) {
      editLoading.value = true
      try {
        const payload: any = { ...editForm.value }
        if (Array.isArray(payload.resolution_tiers)) {
          payload.resolution_tiers = payload.resolution_tiers.join(',')
        }
        if (payload.image_counts !== undefined) {
          payload.image_counts = String(payload.image_counts)
        }
        if (payload.resolution_configs) {
          payload.resolution_configs = JSON.stringify(payload.resolution_configs)
        }

        const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
        const res = await fetch(`${apiUrl}/api/admin/upstreams/update`, {
          method: 'POST',
          headers: getHeaders(),
          body: JSON.stringify(payload)
        })
        const data = await res.json()
        if (res.ok && data.success) {
          ElMessage.success('修改成功')
          showEditDialog.value = false
          fetchUpstreamList()
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
    const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    const res = await fetch(`${apiUrl}/api/admin/upstreams/delete`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ id: row.id })
    })
    const data = await res.json()
    if (res.ok && data.success) {
      ElMessage.success('删除成功')
      fetchUpstreamList()
    } else {
      ElMessage.error(data.error || '删除失败')
    }
  } catch (error) {
    ElMessage.error('网络错误，请稍后重试')
  }
}

const handleSetPrimary = async (row: any) => {
  try {
    const apiUrl = (window as any).APP_CONFIG?.API_BASE_URL || ''
    const res = await fetch(`${apiUrl}/api/admin/upstreams/set-primary`, {
      method: 'POST',
      headers: getHeaders(),
      body: JSON.stringify({ id: row.id, logical_model: row.logical_model, model_type: row.model_type })
    })
    const data = await res.json()
    if (res.ok && data.success) {
      ElMessage.success('设置主用成功')
      fetchUpstreamList()
    } else {
      ElMessage.error(data.error || '设置失败')
    }
  } catch (error) {
    ElMessage.error('网络错误，请稍后重试')
  }
}

onMounted(() => {
  fetchUpstreamList()
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
  padding: 20px;
}

.table-toolbar {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-start;
}

.type-filter-group {
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
  border-radius: 8px;
}

.font-bold {
  font-weight: 600;
}

.text-gray-800 {
  color: #1f2937;
}

.text-gray-500 {
  color: #6b7280;
}

.truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.text-sm {
  font-size: 0.875rem;
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
.add-config-btn-wrapper {
  text-align: center;
  margin: 20px 0;
}

.add-config-btn {
  width: 80%;
  height: 44px;
  border-radius: 8px;
  border: 1px dashed var(--el-color-primary);
  background-color: var(--el-color-primary-light-9);
  color: var(--el-color-primary);
  transition: all 0.3s;
}

.add-config-btn:hover {
  background-color: var(--el-color-primary);
  color: white;
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
  border-style: solid;
}

.add-config-btn .el-icon {
  margin-right: 6px;
  font-size: 16px;
}

.add-config-btn span {
  font-size: 14px;
  font-weight: 500;
  letter-spacing: 1px;
}
</style>
