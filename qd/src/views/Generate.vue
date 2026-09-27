<template>
  <div class="layout-container">
    <!-- Mobile Sidebar Overlay -->
    <div v-if="isMobileMenuOpen" class="sidebar-overlay" @click="toggleMobileMenu"></div>

    <!-- Mini Sidebar -->
<aside class="mini-sidebar" :class="{ 'mobile-open': isMobileMenuOpen }">
  <div class="mini-sidebar-top">
    <div class="logo-mini">
      <img :src="siteLogo" alt="logo" class="logo-img"/>
    </div>
    <div class="nav-items">
      <div class="nav-item" :class="{ active: route.path === '/' }" title="创作" @click="router.push('/')">
        <el-icon><MagicStick /></el-icon>
        <span>创作</span>
      </div>
      <div class="nav-item" :class="{ active: route.path === '/explore' }" title="探索" @click="router.push('/explore'); isMobileMenuOpen = false">
        <el-icon><Compass /></el-icon>
        <span>探索</span>
      </div>
      <div class="nav-item" title="代码" @click="ElMessage.info('开发中'); isMobileMenuOpen = false">
        <el-icon><Monitor /></el-icon>
        <span>代码</span>
      </div>
      <div class="nav-item" title="设计" @click="ElMessage.info('开发中'); isMobileMenuOpen = false">
        <el-icon><Brush /></el-icon>
        <span>设计</span>
      </div>
  </div>
  </div>
</aside>

<!-- History Sidebar -->
<aside v-show="route.path === '/'" class="history-sidebar" :class="{ 'mobile-open': isMobileMenuOpen }">
  <div class="history-header">
    <span class="history-header-title">创作中心</span>
    <el-icon class="new-chat-icon" @click="resetTask" title="新建任务"><Edit /></el-icon>
  </div>
  <div class="sidebar-content" ref="sidebarContentRef" @dragover="handleSidebarDragOver" @dragleave="handleSidebarDragLeave" @drop="stopDragScroll" @dragend="stopDragScroll">
    <!-- 项目组列表 -->
    <div class="project-groups">
      <div class="group-header" @click="isProjectGroupsExpanded = !isProjectGroupsExpanded" style="cursor: pointer; user-select: none; display: flex; align-items: center; justify-content: space-between;">
        <div style="display: flex; align-items: center;">
          <el-icon style="margin-right: 4px; transition: transform 0.3s; font-size: 12px; color: #9ca3af;" :style="{ transform: isProjectGroupsExpanded ? 'rotate(90deg)' : 'rotate(0deg)' }"><ArrowRight /></el-icon>
          <span class="group-title">项目组</span>
          <div v-if="generatingProjectGroupsCount > 0" class="generating-badge" style="display: flex; align-items: center; background-color: #e0f2fe; color: #2563eb; padding: 2px 6px; border-radius: 10px; font-size: 10px; font-weight: bold; margin-left: 8px;">
            <el-icon class="is-loading" style="margin-right: 2px;"><Loading /></el-icon> {{ generatingProjectGroupsCount }}个进行中
          </div>
        </div>
        <el-icon class="add-group-icon" @click.stop="createProjectGroup" title="新建项目组" style="font-size: 32px !important; padding: 8px; width: 32px; height: 32px;"><Plus style="width: 100%; height: 100%;" /></el-icon>
      </div>
      
      <div v-show="isProjectGroupsExpanded">
        <!-- 项目组骨架屏 -->
        <div class="sidebar-skeleton" v-if="isInitialLoading && projectGroups.length === 0">
          <div class="skeleton-history-item" v-for="i in 3" :key="`group-${i}`">
            <div class="skeleton-history-thumbnail" style="border-radius: 4px; width: 24px; height: 24px; background-color: #e5e7eb;"></div>
            <div class="skeleton-history-text" style="height: 14px; max-width: 60%; background-color: #e5e7eb;"></div>
          </div>
        </div>
        
        <!-- 用户创建的项目组 -->
        <div 
          v-else
          v-for="group in projectGroups" 
          :key="group.id"
          class="project-group-container"
        >
        <div
          class="project-group-item"
          :class="{ 'dropdown-open': openDropdownId === group.id }"
          @dragover.prevent
          @drop="handleDropSession($event, group.id)"
          @click="toggleGroup(group.id)"
        >
          <div class="group-folder-icon">
            <el-icon><FolderOpened v-if="!collapsedGroups.has(group.id)" /><Folder v-else /></el-icon>
          </div>
          
          <div class="group-content" v-if="editingGroupId === group.id">
            <el-input
              v-model="editGroupName"
              size="small"
              @input="handleGroupNameInput"
              @blur="saveGroupName(group)"
              @keyup.enter="saveGroupName(group)"
              @click.stop
              ref="editGroupInputRef"
            />
          </div>
          <div v-else style="display: flex; align-items: center; flex: 1; overflow: hidden; min-width: 0;">
            <span class="group-name" :title="group.name" style="margin-right: 8px;">{{ group.name }}</span>
            <div v-if="getGeneratingCountInGroup(group.id) > 0" class="generating-badge" style="display: flex; align-items: center; background-color: #e0f2fe; color: #2563eb; padding: 2px 6px; border-radius: 10px; font-size: 10px; font-weight: bold; flex-shrink: 0; margin-right: 4px;">
              <el-icon class="is-loading" style="margin-right: 2px;"><Loading /></el-icon> 正在图片生成{{ getGeneratingCountInGroup(group.id) }}个
            </div>
          </div>
          
          <div class="group-actions" v-if="editingGroupId !== group.id" @click.stop>
            <el-dropdown trigger="click" @command="(cmd: string) => handleGroupCommand(cmd, group)" @visible-change="(v: boolean) => handleDropdownVisible(v, group.id)">
              <el-icon class="action-icon"><MoreFilled /></el-icon>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="rename"><el-icon><EditPen /></el-icon>重命名</el-dropdown-item>
                  <el-dropdown-item command="delete" class="text-danger"><el-icon><Delete /></el-icon>删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>

        <!-- 组内任务列表 -->
        <div class="group-session-list" v-show="!collapsedGroups.has(group.id)">
          <div 
            v-for="(session, index) in getSessionsByGroupId(group.id)" 
            :key="session.session_id"
            class="history-item"
            :class="{ active: activeSessionId === session.session_id, 'dropdown-open': openDropdownId === session.session_id }"
            @click="selectSession(session.session_id)"
            draggable="true"
            @dragstart="handleDragStartSession($event, session)"
            @dragend="stopDragScroll"
            @dragover.prevent
            @drop.stop="handleDropOnSession($event, session, index, group.id)"
          >
            <el-tooltip
              effect="dark"
              :placement="isMobileMenuOpen ? 'bottom' : 'right'"
              :show-after="300"
            >
              <template #content>
                <div style="line-height: 1.5; max-width: 300px; word-break: break-all;">
                  <div style="font-weight: bold; margin-bottom: 4px; white-space: normal;">{{ session.title || '图片生成任务' }}</div>
                  <div style="color: #ccc; font-size: 12px;">更新于: {{ formatTime(session.updated_at || session.created_at) }}</div>
                </div>
              </template>
              <div class="history-info">
                <div class="history-thumbnail">
                  <el-icon v-if="generatingSessions.includes(session.session_id) || session.status === 0" class="is-loading"><Loading /></el-icon>
                  <img v-else-if="session.image_url" :src="formatImageUrl(session.image_url)" class="thumbnail-img" />
                  <el-icon v-else><Picture /></el-icon>
                </div>
                
                <div class="history-content" v-if="editingSessionId === session.session_id">
                  <el-input
                    v-model="editSessionName"
                    size="small"
                    @blur="saveSessionName(session)"
                    @keyup.enter="saveSessionName(session)"
                    @click.stop
                    ref="editInputRef"
                  />
                </div>
                <span v-else class="history-text">{{ session.title || '图片生成任务' }}</span>
              </div>
            </el-tooltip>
            
            <div class="history-actions" v-if="editingSessionId !== session.session_id" @click.stop>
              <el-dropdown trigger="click" @command="(cmd: string) => handleSessionCommand(cmd, session)" @visible-change="(v: boolean) => handleDropdownVisible(v, session.session_id)">
                <el-icon class="action-icon"><MoreFilled /></el-icon>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="rename"><el-icon><EditPen /></el-icon>重命名</el-dropdown-item>
                    <el-dropdown-item command="moveToGroup"><el-icon><Folder /></el-icon>更换项目组</el-dropdown-item>
                    <el-dropdown-item command="removeFromGroup" class="text-warning"><el-icon><FolderRemove /></el-icon>移除项目组</el-dropdown-item>
                    <el-dropdown-item command="delete" class="text-danger"><el-icon><Delete /></el-icon>删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
          <template v-if="groupPagination[group.id] && groupPagination[group.id].hasMore && getSessionsByGroupId(group.id).length > 0">
            <div v-if="groupPagination[group.id].loadingMore" class="sidebar-skeleton" style="margin-top: 8px;">
              <div class="skeleton-history-item" v-for="i in 3" :key="i">
                <div class="skeleton-history-thumbnail"></div>
                <div class="skeleton-history-text" :style="{ width: `${Math.random() * 40 + 40}%` }"></div>
              </div>
            </div>
            <div v-else class="load-more" @click="fetchSessionsByGroup(group.id, true)">
              <el-icon><ArrowDown /></el-icon>加载更多
            </div>
          </template>
        </div>
      </div>
      </div>
      
      <!-- 项目组分页加载更多 -->
      <template v-if="projectGroupsPagination.hasMore && projectGroups.length > 0 && isProjectGroupsExpanded">
        <div v-if="projectGroupsPagination.loadingMore" class="sidebar-skeleton" style="margin-top: 8px;">
          <div class="skeleton-history-item" v-for="i in 3" :key="`group-skel-${i}`">
            <div class="skeleton-history-thumbnail" style="border-radius: 4px; width: 24px; height: 24px; background-color: #e5e7eb;"></div>
            <div class="skeleton-history-text" style="height: 14px; max-width: 60%; background-color: #e5e7eb;"></div>
          </div>
        </div>
        <div v-else class="load-more" @click="fetchProjectGroups(true)">
          <el-icon><ArrowDown /></el-icon>加载更多项目组
        </div>
      </template>
    </div>

    <!-- 最近任务 -->
    <div 
      class="history-list" 
      style="margin-top: 16px; padding-bottom: 20px;"
      @dragover.prevent
      @drop="handleDropSession($event, 0)"
    >
      <div class="history-title" @click="isRecentTasksExpanded = !isRecentTasksExpanded" style="cursor: pointer; user-select: none; display: flex; align-items: center;">
        <el-icon style="margin-right: 4px; transition: transform 0.3s; font-size: 12px; color: #9ca3af;" :style="{ transform: isRecentTasksExpanded ? 'rotate(90deg)' : 'rotate(0deg)' }"><ArrowRight /></el-icon>
        <span>最近任务</span>
      </div>
      
      <div v-show="isRecentTasksExpanded">
        <!-- 侧边栏骨架屏 -->
        <div class="sidebar-skeleton" v-if="isInitialLoading && sessionList.length === 0">
          <div class="skeleton-history-item" v-for="i in 10" :key="i">
            <div class="skeleton-history-thumbnail"></div>
            <div class="skeleton-history-text" :style="{ width: `${Math.random() * 40 + 40}%` }"></div>
          </div>
        </div>
        
        <div class="group-session-list" v-else>
      <div 
          v-for="(session, index) in getSessionsByGroupId(0)" 
          :key="session.session_id"
          class="history-item"
          :class="{ active: activeSessionId === session.session_id, 'dropdown-open': openDropdownId === session.session_id }"
          @click="selectSession(session.session_id)"
          draggable="true"
          @dragstart="handleDragStartSession($event, session)"
          @dragend="stopDragScroll"
          @dragover.prevent
          @drop.stop="handleDropOnSession($event, session, index, 0)"
        >
          <el-tooltip
            effect="dark"
            :placement="isMobileMenuOpen ? 'bottom' : 'right'"
            :show-after="300"
          >
            <template #content>
              <div style="line-height: 1.5; max-width: 300px; word-break: break-all;">
                <div style="font-weight: bold; margin-bottom: 4px; white-space: normal;">{{ session.title || '图片生成任务' }}</div>
                <div style="color: #ccc; font-size: 12px;">更新于: {{ formatTime(session.updated_at || session.created_at) }}</div>
              </div>
            </template>
            <div class="history-info">
              <div class="history-thumbnail">
                <el-icon v-if="generatingSessions.includes(session.session_id) || session.status === 0" class="is-loading"><Loading /></el-icon>
                <img v-else-if="session.image_url" :src="formatImageUrl(session.image_url)" class="thumbnail-img" />
                <el-icon v-else><Picture /></el-icon>
              </div>
              
              <div class="history-content" v-if="editingSessionId === session.session_id">
                <el-input
                  v-model="editSessionName"
                  size="small"
                  @blur="saveSessionName(session)"
                  @keyup.enter="saveSessionName(session)"
                  @click.stop
                  ref="editInputRef"
                />
              </div>
              <span v-else class="history-text">{{ session.title || '图片生成任务' }}</span>
            </div>
          </el-tooltip>
          
          <div class="history-actions" v-if="editingSessionId !== session.session_id" @click.stop>
            <el-dropdown trigger="click" @command="(cmd: string) => handleSessionCommand(cmd, session)" @visible-change="(v: boolean) => handleDropdownVisible(v, session.session_id)">
              <el-icon class="action-icon"><MoreFilled /></el-icon>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="rename"><el-icon><EditPen /></el-icon>重命名</el-dropdown-item>
                  <el-dropdown-item command="moveToGroup"><el-icon><Folder /></el-icon>移动到项目组</el-dropdown-item>
                  <el-dropdown-item command="delete" class="text-danger"><el-icon><Delete /></el-icon>删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
      <div v-if="getSessionsByGroupId(0).length === 0" class="no-history" style="pointer-events: none;">暂无最近任务</div>
      <template v-if="groupPagination[0] && groupPagination[0].hasMore && getSessionsByGroupId(0).length > 0">
        <div v-if="groupPagination[0].loadingMore" class="sidebar-skeleton" style="margin-top: 8px;">
          <div class="skeleton-history-item" v-for="i in 3" :key="i">
            <div class="skeleton-history-thumbnail"></div>
            <div class="skeleton-history-text" :style="{ width: `${Math.random() * 40 + 40}%` }"></div>
          </div>
        </div>
        <div v-else class="load-more" @click="fetchSessionsByGroup(0, true)">
          <el-icon><ArrowDown /></el-icon>加载更多
        </div>
      </template>
      </div>
    </div>
  </div>

</aside>

<!-- Global Account Container (Spans both sidebars) -->
<div class="global-account-container" :class="{ 'mobile-open': isMobileMenuOpen, 'explore-mode': route.path === '/explore' }">
  <el-dropdown v-if="isLoggedIn" trigger="click" placement="top-start" @command="handleCommand" class="account-dropdown custom-account-dropdown" :teleported="true">
    <button type="button" class="accountTrigger">
      <span class="accountTriggerAvatar">
        <el-avatar :size="32" :icon="UserFilled" />
      </span>
      <div class="accountTriggerInfo" v-show="route.path === '/'">
        <div class="accountTriggerNameWrap">
          <span class="accountTriggerName" :title="username">{{ username }}</span>
          <span class="accountTriggerEmail" :title="userEmail" v-if="userEmail">{{ userEmail }}</span>
        </div>
        <span class="accountTriggerMembership">
          <span class="accountHostTag" style="cursor: pointer;" @click.stop="handleCommand('pointsRecord')" title="查看积分与充值">
            <el-icon style="margin-right: 4px; font-size: 14px; color: #f59e0b;"><Coin /></el-icon>
            {{ parseFloat(Number(userPoints || 0).toFixed(2)) }} 积分
          </span>
        </span>
      </div>
    </button>
    <template #dropdown>
      <el-dropdown-menu class="elegant-profile-dropdown">
        <!-- 用户信息区 -->
        <div class="elegant-header">
          <div class="elegant-avatar-wrap">
            <el-avatar :size="46" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
          </div>
          <div class="elegant-info">
            <div class="elegant-name" :title="username">{{ username }}</div>
            <div class="elegant-email" :title="userEmail" v-if="userEmail">{{ userEmail }}</div>
          </div>
        </div>

        <!-- 积分卡片区 -->
        <div class="elegant-points-card">
          <div class="elegant-points-left">
            <span class="elegant-points-label">我的可用积分</span>
            <div class="elegant-points-value">
              <el-icon><Coin /></el-icon>
              <span>{{ parseFloat(Number(userPoints || 0).toFixed(2)) }}</span>
            </div>
          </div>
          <div class="elegant-points-right">
            <el-button round size="small" class="elegant-recharge-btn" @click="handleCommand('pointsRecord')">获取积分</el-button>
          </div>
        </div>

        <div class="elegant-divider"></div>

        <!-- 菜单操作区 -->
        <div class="elegant-menu-group">
          <el-dropdown-item command="pointsRecord" class="elegant-menu-item">
            <el-icon class="elegant-icon"><List /></el-icon>
            <span>积分明细</span>
            <el-icon class="elegant-arrow"><ArrowRight /></el-icon>
          </el-dropdown-item>
          
          <el-dropdown-item command="settings" class="elegant-menu-item">
            <el-icon class="elegant-icon"><Setting /></el-icon>
            <span>个人设置</span>
            <el-icon class="elegant-arrow"><ArrowRight /></el-icon>
          </el-dropdown-item>
        </div>

        <div class="elegant-divider"></div>

        <div class="elegant-menu-group">
          <el-dropdown-item command="logout" class="elegant-menu-item elegant-logout">
            <el-icon class="elegant-icon"><SwitchButton /></el-icon>
            <span>退出账号</span>
          </el-dropdown-item>
        </div>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
  
  <button v-else type="button" class="accountTrigger" @click="authStore.openLogin()" title="点击登录">
    <span class="accountTriggerAvatar">
      <el-avatar :size="32" :icon="UserFilled" />
    </span>
    <div class="accountTriggerInfo" v-show="route.path === '/'">
      <div class="accountTriggerNameWrap">
        <span class="accountTriggerName">未登录</span>
      </div>
      <span class="accountTriggerMembership">
        <span class="accountHostTag login-tag">点击登录</span>
      </span>
    </div>
  </button>
</div>

<!-- Main Content -->
<main class="main-container" v-show="route.path === '/'">
  <!-- Header -->
  <header class="main-header">
    <div class="header-left">
      <div class="mobile-menu-btn" @click="toggleMobileMenu">
        <el-icon><Expand /></el-icon>
      </div>
      <div class="header-title-container">
        <el-icon class="header-title-icon"><Document /></el-icon>
        <div class="header-title-content">
          <span v-if="!isNewTaskMode && currentGroupName" class="header-group-name" v-cloak>
            {{ currentGroupName }} <span class="header-separator">/</span>
          </span>
          <span class="header-title" v-cloak>
            {{ isNewTaskMode ? '新建图片生成' : (currentSessionTitle || '图片生成任务') }}
          </span>
        </div>
      </div>
    </div>
    <div class="header-actions">
      <div v-if="!isLoggedIn" class="header-login-btn" @click="authStore.openLogin()">
        登录 / 注册
      </div>
      <div v-else class="header-points-btn" @click="handleCommand('pointsRecord')" title="查看积分记录">
        <el-icon><Coin /></el-icon>
        <span>{{ parseFloat(Number(userPoints || 0).toFixed(2)) }} 积分</span>
      </div>
    </div>
  </header>

  <!-- Chat / Result Area -->
  <div class="content-area" ref="contentAreaRef" @scroll="handleContentScroll">
    <!-- Welcome Screen -->
    <div v-if="isNewTaskMode || (!isInitialLoading && sessionList.length === 0 && (!activeSessionId || !tempTasksMap[activeSessionId] || tempTasksMap[activeSessionId].length === 0))" class="welcome-screen">
      <div class="welcome-logo">
        <el-icon :size="48" color="#333"><MagicStick /></el-icon>
      </div>
      <h1 class="welcome-title">使用 AI 图片生成工作</h1>
      <p class="welcome-desc">
        帮助你将想法转化为精美的图片，输出专业的可视化结果。
      </p>
      <div class="suggestion-cards">
        <div class="suggestion-card" @click="useSuggestion('一只可爱的猫咪在太空遨游')">
          <el-icon><Sunny /></el-icon>
          <span>一只可爱的猫咪在太空遨游</span>
        </div>
        <div class="suggestion-card" @click="useSuggestion('赛博朋克风格的未来城市，霓虹灯闪烁')">
          <el-icon><Star /></el-icon>
          <span>赛博朋克风格的未来城市</span>
        </div>
      </div>
    </div>

    <!-- Empty Session Screen -->
    <div v-else-if="!isTasksLoading && (!activeSessionId || !tempTasksMap[activeSessionId] || tempTasksMap[activeSessionId].length === 0) && chatList.length === 0" class="empty-session-screen">
      <el-empty description="当前会话内容已清空，快来探索生成新图片吧！">
        <template #image>
          <div class="custom-empty-icon">
            <div class="icon-layer layer-1"></div>
            <div class="icon-layer layer-2"></div>
            <div class="icon-layer layer-3">
              <el-icon><MagicStick /></el-icon>
            </div>
          </div>
        </template>
      </el-empty>
    </div>

    <!-- Result Screen (Chat Flow) -->
    <div v-else class="result-screen">
      <div v-if="isTasksLoading && currentTaskPage === 1" class="initial-tasks-loading">
        <!-- Skeleton Group 1 -->
        <div class="skeleton-group" style="opacity: 0.2;">
          <div class="message user-message skeleton-message">
            <div class="message-inner" style="width: 100%;">
              <div class="message-body" style="display: flex; flex-direction: column; align-items: flex-end; width: 100%;">
                <div class="user-name skeleton-bg" style="width: 60px; height: 16px; margin-bottom: 2px; margin-right: 2px; border-radius: 4px;"></div>
                <div class="message-content user-content skeleton-bg user-text-skeleton"></div>
              </div>
              <div class="message-avatar">
                <div class="skeleton-avatar-circle"></div>
              </div>
            </div>
          </div>
          <div class="message ai-message skeleton-message">
            <div class="message-inner">
              <div class="message-avatar">
                <div class="skeleton-avatar-circle"></div>
              </div>
              <div class="message-body" style="display: flex; flex-direction: column; align-items: flex-start; width: 100%;">
                <div class="ai-model-name skeleton-bg" style="width: 80px; height: 16px; margin-bottom: 2px; margin-left: 2px; border-radius: 4px;"></div>
                <div class="message-content ai-content skeleton-bg image-skeleton small"></div>
              </div>
            </div>
          </div>
        </div>
        <!-- Skeleton Group 2 -->
        <div class="skeleton-group" style="opacity: 0.4;">
          <div class="message user-message skeleton-message">
            <div class="message-inner" style="width: 100%;">
              <div class="message-body" style="display: flex; flex-direction: column; align-items: flex-end; width: 100%;">
                <div class="user-name skeleton-bg" style="width: 60px; height: 16px; margin-bottom: 2px; margin-right: 2px; border-radius: 4px;"></div>
                <div class="message-content user-content skeleton-bg user-text-skeleton medium"></div>
              </div>
              <div class="message-avatar">
                <div class="skeleton-avatar-circle"></div>
              </div>
            </div>
          </div>
          <div class="message ai-message skeleton-message">
            <div class="message-inner">
              <div class="message-avatar">
                <div class="skeleton-avatar-circle"></div>
              </div>
              <div class="message-body" style="display: flex; flex-direction: column; align-items: flex-start; width: 100%;">
                <div class="ai-model-name skeleton-bg" style="width: 80px; height: 16px; margin-bottom: 2px; margin-left: 2px; border-radius: 4px;"></div>
                <div class="message-content ai-content skeleton-bg image-skeleton large"></div>
              </div>
            </div>
          </div>
        </div>
        <!-- Skeleton Group 3 -->
        <div class="skeleton-group" style="opacity: 0.7;">
          <div class="message user-message skeleton-message">
            <div class="message-inner" style="width: 100%;">
              <div class="message-body" style="display: flex; flex-direction: column; align-items: flex-end; width: 100%;">
                <div class="user-name skeleton-bg" style="width: 60px; height: 16px; margin-bottom: 2px; margin-right: 2px; border-radius: 4px;"></div>
                <div class="message-content user-content skeleton-bg user-text-skeleton short"></div>
              </div>
              <div class="message-avatar">
                <div class="skeleton-avatar-circle"></div>
              </div>
            </div>
          </div>
          <div class="message ai-message skeleton-message">
            <div class="message-inner">
              <div class="message-avatar">
                <div class="skeleton-avatar-circle"></div>
              </div>
              <div class="message-body" style="display: flex; flex-direction: column; align-items: flex-start; width: 100%;">
                <div class="ai-model-name skeleton-bg" style="width: 80px; height: 16px; margin-bottom: 2px; margin-left: 2px; border-radius: 4px;"></div>
                <div class="message-content ai-content skeleton-bg image-skeleton small"></div>
              </div>
            </div>
          </div>
        </div>
        <!-- Skeleton Group 4 -->
        <div class="skeleton-group" style="opacity: 1;">
          <div class="message user-message skeleton-message">
            <div class="message-inner" style="width: 100%;">
              <div class="message-body" style="display: flex; flex-direction: column; align-items: flex-end; width: 100%;">
                <div class="user-name skeleton-bg" style="width: 60px; height: 16px; margin-bottom: 2px; margin-right: 2px; border-radius: 4px;"></div>
                <div class="message-content user-content skeleton-bg user-text-skeleton"></div>
              </div>
              <div class="message-avatar">
                <div class="skeleton-avatar-circle"></div>
              </div>
            </div>
          </div>
          <div class="message ai-message skeleton-message">
            <div class="message-inner">
              <div class="message-avatar">
                <div class="skeleton-avatar-circle"></div>
              </div>
              <div class="message-body" style="display: flex; flex-direction: column; align-items: flex-start; width: 100%;">
                <div class="ai-model-name skeleton-bg" style="width: 80px; height: 16px; margin-bottom: 2px; margin-left: 2px; border-radius: 4px;"></div>
                <div class="message-content ai-content skeleton-bg image-skeleton medium"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <template v-else>
        <div class="load-more-chat" v-if="hasMoreTasks" @click="!isLoadingMoreTasks && fetchTasksBySession(activeSessionId!, true)">
          <div class="load-more-inner" v-if="!isLoadingMoreTasks">
            <el-icon class="mr-1"><Clock /></el-icon>
            <span>加载更多历史记录...</span>
          </div>
          <div class="load-more-inner" v-else>
            <el-icon class="mr-1 is-loading"><Loading /></el-icon>
            <span>加载中...</span>
          </div>
        </div>

        <template v-for="(task, index) in chatList" :key="task.id">
        <!-- Chat Time Divider -->
        <div class="chat-time-divider" v-if="task.created_at && shouldShowChatTime(task, index)">
          <span>{{ formatChatTime(task.created_at) }}</span>
        </div>
        <!-- User Message -->
        <div class="message user-message" :id="'task-' + task.id">
          <div class="message-inner">
            <div class="message-body" style="display: flex; flex-direction: column; align-items: flex-end;">
              <div class="user-name" style="font-size: 13px; color: #6b7280; margin-bottom: 10px; margin-right: 4px; font-weight: 500;">{{ username }}</div>
              <div class="message-content user-content">
                <div v-if="task.reference_image" class="user-reference-image">
                  <template v-for="(img, idx) in parseImageUrls(task.reference_image)" :key="idx">
                    <div class="reference-group" :style="{ zIndex: parseImageUrls(task.reference_image).length - idx }">
                      <div class="reference-item" :style="{ transform: `rotate(${[-8, 5, -8, 5, -8, 4][idx % 6]}deg)` }">
                        <el-image
                          :src="formatImageUrl(img)"
                          :preview-src-list="parseImageUrls(task.reference_image).map(i => formatImageUrl(i))"
                          :initial-index="idx"
                          :preview-teleported="true"
                          :hide-on-click-modal="true"
                          fit="cover"
                          class="ref-img-display"
                          @load="scrollToBottom(false)"
                        />
                      </div>
                    </div>
                  </template>
                </div>
                {{ task.prompt }}
              </div>
              <div class="message-meta" :class="{ 'is-active': activePopoverId === task.id }">
                <span class="meta-date" v-if="task.created_at">{{ formatTime(task.created_at) }}</span>
                <span class="meta-action-btn" @click="copyPrompt(task.prompt)" title="复制提示词"><el-icon><CopyDocument /></el-icon>复制</span>
                  <el-popover
                    placement="bottom-end"
                    :width="240"
                    trigger="hover"
                    popper-class="custom-detail-popover"
                    :show-after="200"
                    @show="activePopoverId = task.id"
                    @hide="activePopoverId = null"
                  >
                    <template #reference>
                      <span class="simple-detail-btn"><el-icon><Operation /></el-icon>生图参数</span>
                    </template>
                  <div class="task-detail-content">
                    <div class="detail-header">
                      <el-icon><Document /></el-icon>
                      <span>生成参数详情</span>
                    </div>
                    <div class="detail-body">
                      <div class="detail-row" v-if="task.series_id">
                        <span class="label"><el-icon><Cpu /></el-icon>模型</span> 
                        <span class="value">{{ getModelName(task.series_id) }}</span>
                      </div>
                      <div class="detail-row" v-if="task.resolution">
                        <span class="label"><el-icon><Monitor /></el-icon>分辨率</span> 
                        <span class="value">{{ task.resolution }}</span>
                      </div>
                      <div class="detail-row">
                        <span class="label"><el-icon><Crop /></el-icon>比例</span> 
                        <span class="value">{{ task.size || 'auto' }}</span>
                      </div>
                      <div class="detail-row">
                        <span class="label"><el-icon><CopyDocument /></el-icon>张数</span> 
                        <span class="value">{{ task.num_images || 1 }} 张</span>
                      </div>
                      <div class="detail-row" v-if="task.created_at">
                        <span class="label"><el-icon><Clock /></el-icon>时间</span> 
                        <span class="value">{{ formatTime(task.created_at) }}</span>
                      </div>
                      <div class="detail-row highlight-row">
                        <span class="label"><el-icon><Coin /></el-icon>消耗积分</span> 
                        <span class="value highlight-value">-{{ parseFloat((getModelCredits(task.series_id, task.resolution) * (task.num_images || 1)).toFixed(2)) }}</span>
                      </div>
                      <div class="detail-row highlight-row" v-if="task.status === 2 || task.status === 3">
                        <span class="label"><el-icon><RefreshLeft /></el-icon>退还积分</span> 
                        <span class="value" style="color: #67c23a; font-weight: bold;">+{{ calculateRefundedPoints(task) }}</span>
                      </div>
                    </div>
                  </div>
                </el-popover>
              </div>
            </div>
            <div class="message-avatar">
              <el-avatar :size="36" :icon="UserFilled" />
            </div>
          </div>
        </div>

        <!-- AI Message -->
        <div class="message ai-message">
          <div class="message-inner">
            <div class="message-avatar">
              <div class="ai-avatar">
                <img :src="siteLogo" alt="AI" style="width: 36px; height: 36px; object-fit: cover; border-radius: 50%;" />
              </div>
            </div>
            <div class="message-body" style="display: flex; flex-direction: column; align-items: flex-start; width: 100%; min-width: 0;">
              <div class="ai-model-name" style="font-size: 13px; color: #6b7280; margin-bottom: 10px; margin-left: 4px; font-weight: 500;">{{ getModelName(task.series_id) }}</div>
              <div class="message-content ai-content">
                <div v-if="task.status === 0" class="image-result-grid">
                  <div :class="['image-grid-container', `grid-count-${task.num_images || 1}`]">
                    <!-- Render generated partial images -->
                    <template v-for="(url, idx) in parseImageUrls(task.image_url)" :key="`partial-${idx}`">
                      <div class="grid-item" :style="getAspectRatioStyle(task.size)">
                        <el-image
                          :src="formatImageUrl(url)"
                          :preview-src-list="parseImageUrls(task.image_url).map(img => formatImageUrl(img))"
                          :initial-index="idx"
                          :preview-teleported="true"
                          :hide-on-click-modal="true"
                          fit="contain"
                          class="generated-image"
                          @load="scrollToBottom(false)"
                        >
                          <template #error>
                            <div class="image-slot">
                              <el-icon><Picture /></el-icon>
                              <span style="font-size: 14px; margin-top: 10px;">加载失败</span>
                            </div>
                          </template>
                        </el-image>
                        <div class="image-actions">
                          <el-button type="primary" link :icon="Download" @click="downloadImage(formatImageUrl(url))">下载</el-button>
                        </div>
                      </div>
                    </template>
                    
                    <!-- Render loading placeholders for the rest -->
                    <div v-for="n in Math.max(0, (task.num_images || 1) - parseImageUrls(task.image_url).length)" :key="`loading-${n}`" class="grid-item" :style="getAspectRatioStyle(task.size)">
                      <div class="generating-placeholder">
                        <div class="placeholder-bg"></div>
                        <div class="placeholder-scan"></div>
                        <div class="placeholder-content">
                          <el-icon class="generating-spinner is-loading"><Loading /></el-icon>
                          <span class="generating-text">AI 渲染中...</span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <div v-else-if="task.status === 1 || task.status === 3" class="image-result-grid">
                <div :class="['image-grid-container', `grid-count-${task.num_images || 1}`]" v-if="task.image_url && task.image_url !== '[]'">
                  <template v-for="(url, idx) in parseImageUrls(task.image_url)" :key="idx">
                    <div class="grid-item" :style="getAspectRatioStyle(task.size)">
                      <el-image
                        :src="formatImageUrl(url)"
                        :preview-src-list="parseImageUrls(task.image_url).map(img => formatImageUrl(img))"
                        :initial-index="idx"
                        :preview-teleported="true"
                        :hide-on-click-modal="true"
                        fit="contain"
                        class="generated-image"
                        @load="scrollToBottom(false)"
                      >
                        <template #error>
                          <div class="image-slot">
                            <el-icon><Picture /></el-icon>
                            <span style="font-size: 14px; margin-top: 10px;">加载失败</span>
                          </div>
                        </template>
                      </el-image>
                      <div class="image-actions">
                        <el-button type="primary" link :icon="Download" @click="downloadImage(formatImageUrl(url))">下载</el-button>
                      </div>
                    </div>
                  </template>
                  
                  <!-- Render failed placeholders for partial success -->
                  <div v-for="n in Math.max(0, (task.num_images || 1) - parseImageUrls(task.image_url).length)" :key="`failed-${n}`" class="grid-item" style="display: flex; flex-direction: column; width: 100%; height: 100%;">
                    <div class="image-slot" :style="getAspectRatioStyle(task.size)" style="background-color: #fef2f2; border: 1px dashed #fca5a5; display: flex; flex-direction: column; justify-content: center; align-items: center; border-radius: 8px; width: 100%; height: 100%; flex: 1; min-height: 120px; box-sizing: border-box; padding: 8px; text-align: center;">
                      <el-icon color="#f87171" :size="20" style="flex-shrink: 0; margin-bottom: 4px;"><CircleClose /></el-icon>
                      <span style="font-size: 11px; color: #ef4444; word-break: break-word; white-space: normal; line-height: 1.2; max-width: 100%;">生成失败已退款<br/><br/><span style="opacity: 0.8;">{{ getErrorMessage(task.log_content, task.status) }}</span></span>
                    </div>
                  </div>
                </div>
                <div v-else class="image-slot" style="min-height: 100px;">
                  <el-icon color="#e6a23c"><Warning /></el-icon>
                  <span style="font-size: 14px; margin-top: 10px; color: #e6a23c; white-space: pre-wrap;">
                    生成成功，但未返回图片数据
                  </span>
                </div>
              </div>
              <div v-else class="image-slot" style="min-height: 100px;">
                  <el-icon color="#f56c6c"><Warning /></el-icon>
                  <span style="font-size: 14px; margin-top: 10px; color: #f56c6c; white-space: pre-wrap; word-break: break-all;">
    {{ getErrorMessage(task.log_content, task.status) }}
  </span>
                </div>
              </div>
              <div class="message-meta" style="justify-content: flex-start; margin-left: 17px; margin-top: 8px;">
                <span v-if="task.status !== 0" class="meta-action-btn always-show" @click="regenerateTask(task)" title="重新生成">
                  <el-icon><RefreshRight /></el-icon>重新生成
                </span>
                <span v-if="task.status !== 0" class="meta-action-btn always-show" @click="reEditTask(task)" title="重新编辑">
                  <el-icon><EditPen /></el-icon>重新编辑
                </span>
                <el-dropdown v-if="task.status !== 0" trigger="click" @command="(cmd: string) => handleTaskCommand(cmd, task)">
                  <span class="meta-action-btn always-show" style="margin-right: 0;" title="更多">
                    <el-icon><MoreFilled /></el-icon>
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item command="delete" style="color: #ef4444;">
                        <el-icon><Delete /></el-icon>删除该批次结果
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
            </div>
          </div>
        </div>
      </template>
      </template>
    </div>
  </div>

    <!-- Input Area -->
  <div class="input-area-wrapper" :class="{ 'is-shrunk': isInputShrunk }">
    <div class="input-tools-container mobile-only">
      <div class="input-tools">
        
        <el-dropdown class="model-dropdown-wrapper" trigger="click" placement="top-start" @command="(val: string) => form.series_id = val">
          <button class="combined-settings-btn model-settings-btn" type="button">
            <div class="model-btn-content">
              <el-icon class="cpu-icon"><Cpu /></el-icon>
              <span class="btn-text model-name-text" :title="availableModels.find(m => m.series_id === form.series_id)?.name || '默认模型'">
                {{ availableModels.find(m => m.series_id === form.series_id)?.name || '默认模型' }}
              </span>
              <span v-if="availableModels.find(m => m.series_id === form.series_id)?.activity_tag" class="shimmer-tag" :style="{ fontSize: '10px', color: '#fff', background: availableModels.find(m => m.series_id === form.series_id)?.activity_tag_color || '#10b981', padding: '2px 4px', borderRadius: '4px', marginLeft: '6px', whiteSpace: 'nowrap', flexShrink: 0 }">
                {{ availableModels.find(m => m.series_id === form.series_id)?.activity_tag }}
              </span>
              <span v-else-if="availableModels.find(m => m.series_id === form.series_id) && hasFreeResolution(availableModels.find(m => m.series_id === form.series_id))" class="shimmer-tag" style="font-size: 10px; color: #fff; background: #10b981; padding: 2px 4px; border-radius: 4px; margin-left: 6px; white-space: nowrap; flex-shrink: 0;">
                限时免费
              </span>
            </div>
            <el-icon class="arrow-icon"><ArrowDown /></el-icon>
          </button>
          <template #dropdown>
            <el-dropdown-menu class="model-dropdown-menu">
              <el-dropdown-item 
                v-for="model in availableModels" 
                :key="model.series_id" 
                :command="model.series_id"
                :class="{ 'is-active-model': form.series_id === model.series_id }"
              >
                <div style="display: flex; align-items: center; justify-content: space-between; width: 100%; gap: 12px;">
                    <span>{{ model.name }}</span>
                    <span v-if="model.activity_tag" class="shimmer-tag" :style="{ fontSize: '11px', color: '#fff', background: model.activity_tag_color || '#10b981', padding: '2px 6px', borderRadius: '4px', lineHeight: '1.2', whiteSpace: 'nowrap' }">{{ model.activity_tag }}</span>
                    <span v-else-if="hasFreeResolution(model)" class="shimmer-tag" style="font-size: 11px; color: #fff; background: #10b981; padding: 2px 6px; border-radius: 4px; line-height: 1.2; white-space: nowrap;">限时免费</span>
                  </div>
              </el-dropdown-item>
              <el-dropdown-item 
                v-if="availableModels.length === 0" 
                command="default"
                :class="{ 'is-active-model': form.series_id === 'default' }"
              >
                默认模型
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        
        <el-popover ref="desktopPopoverRef" placement="top-start" :width="windowWidth <= 768 ? windowWidth - 32 : 360" trigger="click" popper-class="settings-popover" :popper-options="{ modifiers: [{ name: 'preventOverflow', options: { padding: 16 } }] }">
          <template #reference>
            <button class="combined-settings-btn" type="button">
              <el-icon><Crop /></el-icon>
              <span class="btn-text">
                <template v-if="currentModelRatios.length > 0">
                  {{ form.aspect_ratio || 'auto' }}
                  <div class="divider"></div>
                </template>
                {{ form.resolution || '1K' }}
                <template v-if="currentModelImageCounts.length > 0">
                  <div class="divider"></div>
                  <span>{{ form.num_images || 1 }}<span class="unit-text">张</span></span>
                </template>
              </span>
            </button>
          </template>
          <div class="settings-panel">
            <div class="setting-item" v-if="currentModelRatios.length > 0">
              <div class="setting-label">比例</div>
              <div class="setting-options">
                <div class="option-btn" v-for="ratio in currentModelRatios" :key="ratio" :class="{ active: form.aspect_ratio === ratio }" @click="form.aspect_ratio = ratio">
                  {{ ratio }}
                </div>
              </div>
            </div>
            <div class="setting-item">
              <div class="setting-label">分辨率</div>
              <div class="setting-options">
                <div class="option-btn" v-for="res in currentModelResolutions" :key="res" :class="{ active: form.resolution === res }" @click="form.resolution = res">
                  <div style="display: flex; align-items: center; gap: 6px;">
                    <span>{{ res }}</span>
                    <span v-if="getTierConfigForRes(currentModelInfo, res, 'enabled', true) === 'maintenance'" class="shimmer-tag" style="font-size: 11px; color: #fff; background: #9ca3af; padding: 0 4px; border-radius: 4px; line-height: 1.4; white-space: nowrap;">维护中</span>
                    <span v-else-if="getTierConfigForRes(currentModelInfo, res, 'tag', '')" class="shimmer-tag" :style="{ fontSize: '11px', color: '#fff', background: getTierConfigForRes(currentModelInfo, res, 'tag_color', '#10b981'), padding: '0 4px', borderRadius: '4px', lineHeight: '1.4', whiteSpace: 'nowrap' }">{{ getTierConfigForRes(currentModelInfo, res, 'tag', '') }}</span>
                    <span v-else-if="Number(getTierConfigForRes(currentModelInfo, res, 'credits_per_image', 1)) === 0" class="shimmer-tag" style="font-size: 11px; color: #fff; background: #10b981; padding: 0 4px; border-radius: 4px; line-height: 1.4; white-space: nowrap;">限时免费</span>
                    <span v-else style="font-size: 11px; opacity: 0.7;">{{ parseFloat(Number(getTierConfigForRes(currentModelInfo, res, 'credits_per_image', 1)).toFixed(2)) }}积分</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="setting-item" v-if="currentModelImageCounts.length > 0">
              <div class="setting-label">生成数量</div>
              <div class="setting-options">
                <div class="option-btn" v-for="count in currentModelImageCounts" :key="count" :class="{ active: form.num_images === Number(count) }" @click="form.num_images = Number(count)">
                  {{ count }}张
                </div>
              </div>
            </div>
          </div>
        </el-popover>
      </div>
    </div>
    <div class="input-container">
      <!-- Reference Image Preview -->
      <div class="reference-image-preview" v-if="false">
        <div class="preview-list">
          <div class="preview-container" v-for="(img, index) in form.reference_image" :key="index" :style="{ transform: `rotate(${index === 0 ? 8 : (index === 1 ? -4 : 22)}deg)`, zIndex: index, left: `${index * 24}px` }">
            <img :src="img" class="preview-img" />
            <div class="remove-btn" @click.stop="clearReferenceImage(index)">
              <el-icon><Close /></el-icon>
            </div>
          </div>
          <div class="reference-upload-item" @click.stop="triggerUpload" v-if="form.reference_image.length < currentModelMaxRefImages" :style="{ transform: `rotate(${form.reference_image.length === 0 ? 8 : (form.reference_image.length === 1 ? -4 : 22)}deg)`, zIndex: form.reference_image.length, left: `${form.reference_image.length * 24}px` }">
            <div class="reference-upload-content">
              <svg width="1em" height="1em" viewBox="0 0 24 24" preserveAspectRatio="xMidYMid meet" fill="none" role="presentation" xmlns="http://www.w3.org/2000/svg" class="upload-icon-svg"><g><path data-follow-fill="currentColor" d="M10.8 20a1.2 1.2 0 0 0 2.4 0v-6.8H20a1.2 1.2 0 1 0 0-2.4h-6.8V4a1.2 1.2 0 0 0-2.4 0v6.8H4a1.2 1.2 0 0 0 0 2.4h6.8V20Z" clip-rule="evenodd" fill-rule="evenodd" fill="currentColor"></path></g></svg>
            </div>
          </div>
        </div>
      </div>
      
      <div class="input-tools desktop-only">
        
        <el-dropdown class="model-dropdown-wrapper" trigger="click" placement="top-start" @command="(val: string) => form.series_id = val">
          <button class="combined-settings-btn model-settings-btn" type="button">
            <div class="model-btn-content">
              <el-icon class="cpu-icon"><Cpu /></el-icon>
              <span class="btn-text model-name-text" :title="availableModels.find(m => m.series_id === form.series_id)?.name || '默认模型'">
                {{ availableModels.find(m => m.series_id === form.series_id)?.name || '默认模型' }}
              </span>
              <span v-if="availableModels.find(m => m.series_id === form.series_id)?.activity_tag" class="shimmer-tag" :style="{ fontSize: '10px', color: '#fff', background: availableModels.find(m => m.series_id === form.series_id)?.activity_tag_color || '#10b981', padding: '2px 4px', borderRadius: '4px', marginLeft: '6px', whiteSpace: 'nowrap', flexShrink: 0 }">
                {{ availableModels.find(m => m.series_id === form.series_id)?.activity_tag }}
              </span>
              <span v-else-if="availableModels.find(m => m.series_id === form.series_id) && hasFreeResolution(availableModels.find(m => m.series_id === form.series_id))" class="shimmer-tag" style="font-size: 10px; color: #fff; background: #10b981; padding: 2px 4px; border-radius: 4px; margin-left: 6px; white-space: nowrap; flex-shrink: 0;">
                限时免费
              </span>
            </div>
            <el-icon class="arrow-icon"><ArrowDown /></el-icon>
          </button>
          <template #dropdown>
            <el-dropdown-menu class="model-dropdown-menu">
              <el-dropdown-item 
                v-for="model in availableModels" 
                :key="model.series_id" 
                :command="model.series_id"
                :class="{ 'is-active-model': form.series_id === model.series_id }"
              >
                <div style="display: flex; align-items: center; justify-content: space-between; width: 100%; gap: 12px;">
                    <span>{{ model.name }}</span>
                    <span v-if="model.activity_tag" class="shimmer-tag" :style="{ fontSize: '11px', color: '#fff', background: model.activity_tag_color || '#10b981', padding: '2px 6px', borderRadius: '4px', lineHeight: '1.2', whiteSpace: 'nowrap' }">{{ model.activity_tag }}</span>
                    <span v-else-if="hasFreeResolution(model)" class="shimmer-tag" style="font-size: 11px; color: #fff; background: #10b981; padding: 2px 6px; border-radius: 4px; line-height: 1.2; white-space: nowrap;">限时免费</span>
                  </div>
              </el-dropdown-item>
              <el-dropdown-item 
                v-if="availableModels.length === 0" 
                command="default"
                :class="{ 'is-active-model': form.series_id === 'default' }"
              >
                默认模型
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        
        <el-popover ref="mobilePopoverRef" placement="top-start" :width="windowWidth <= 768 ? windowWidth - 32 : 360" trigger="click" popper-class="settings-popover" :popper-options="{ modifiers: [{ name: 'preventOverflow', options: { padding: 16 } }] }">
          <template #reference>
            <button class="combined-settings-btn" type="button">
              <el-icon><Crop /></el-icon>
              <span class="btn-text">
                <template v-if="currentModelRatios.length > 0">
                  {{ form.aspect_ratio || 'auto' }}
                  <div class="divider"></div>
                </template>
                {{ form.resolution || '1K' }}
                <template v-if="currentModelImageCounts.length > 0">
                  <div class="divider"></div>
                  <span>{{ form.num_images || 1 }}<span class="unit-text">张</span></span>
                </template>
              </span>
            </button>
          </template>
          <div class="settings-panel">
            <div class="setting-item" v-if="currentModelRatios.length > 0">
              <div class="setting-label">比例</div>
              <div class="setting-options">
                <div class="option-btn" v-for="ratio in currentModelRatios" :key="ratio" :class="{ active: form.aspect_ratio === ratio }" @click="form.aspect_ratio = ratio">
                  {{ ratio }}
                </div>
              </div>
            </div>
            <div class="setting-item">
              <div class="setting-label">分辨率</div>
              <div class="setting-options">
                <div class="option-btn" v-for="res in currentModelResolutions" :key="res" :class="{ active: form.resolution === res }" @click="form.resolution = res">
                  <div style="display: flex; align-items: center; gap: 6px;">
                    <span>{{ res }}</span>
                    <span v-if="getTierConfigForRes(currentModelInfo, res, 'enabled', true) === 'maintenance'" class="shimmer-tag" style="font-size: 11px; color: #fff; background: #9ca3af; padding: 0 4px; border-radius: 4px; line-height: 1.4; white-space: nowrap;">维护中</span>
                    <span v-else-if="getTierConfigForRes(currentModelInfo, res, 'tag', '')" class="shimmer-tag" :style="{ fontSize: '11px', color: '#fff', background: getTierConfigForRes(currentModelInfo, res, 'tag_color', '#10b981'), padding: '0 4px', borderRadius: '4px', lineHeight: '1.4', whiteSpace: 'nowrap' }">{{ getTierConfigForRes(currentModelInfo, res, 'tag', '') }}</span>
                    <span v-else-if="Number(getTierConfigForRes(currentModelInfo, res, 'credits_per_image', 1)) === 0" class="shimmer-tag" style="font-size: 11px; color: #fff; background: #10b981; padding: 0 4px; border-radius: 4px; line-height: 1.4; white-space: nowrap;">限时免费</span>
                    <span v-else style="font-size: 11px; opacity: 0.7;">{{ parseFloat(Number(getTierConfigForRes(currentModelInfo, res, 'credits_per_image', 1)).toFixed(2)) }}积分</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="setting-item" v-if="currentModelImageCounts.length > 0">
              <div class="setting-label">生成数量</div>
              <div class="setting-options">
                <div class="option-btn" v-for="count in currentModelImageCounts" :key="count" :class="{ active: form.num_images === Number(count) }" @click="form.num_images = Number(count)">
                  {{ count }}张
                </div>
              </div>
            </div>
          </div>
        </el-popover>
      </div>
      
      <div class="input-box" :class="{ 'has-refs': form.reference_image && form.reference_image.length > 0, 'is-focused': isInputFocused }" @click="isInputShrunk ? (isInputFocused = true, isInputShrunk = false) : null">
        <div class="reference-upload-item initial-upload-item" @click.stop="isInputShrunk ? (isInputFocused = true, isInputShrunk = false) : triggerUpload()" v-if="(!form.reference_image || form.reference_image.length === 0) && currentModelMaxRefImages > 0">
          <div class="reference-upload-content" style="transform: rotate(8deg);">
            <svg width="1em" height="1em" viewBox="0 0 24 24" preserveAspectRatio="xMidYMid meet" fill="none" role="presentation" xmlns="http://www.w3.org/2000/svg" class="upload-icon-svg"><g><path data-follow-fill="currentColor" d="M10.8 20a1.2 1.2 0 0 0 2.4 0v-6.8H20a1.2 1.2 0 1 0 0-2.4h-6.8V4a1.2 1.2 0 0 0-2.4 0v6.8H4a1.2 1.2 0 0 0 0 2.4h6.8V20Z" clip-rule="evenodd" fill-rule="evenodd" fill="currentColor"></path></g></svg>
          </div>
        </div>
        
        <!-- 叠加卡片式预览区（移入输入框内部左侧，当有图片时显示） -->
      <div class="reference-image-preview-inline" v-if="form.reference_image && form.reference_image.length > 0" :style="{ '--total-items': (form.reference_image.length + (form.reference_image.length < currentModelMaxRefImages ? 1 : 0)) }">
        <div class="preview-list-inline" :class="{ 'is-expandable': form.reference_image.length > 1 || (form.reference_image.length === 1 && currentModelMaxRefImages > 1), 'is-mobile-expanded': isMobileRefExpanded }" @click="isInputShrunk ? (isInputFocused = true, isInputShrunk = false) : (windowWidth <= 768 ? isMobileRefExpanded = true : null)">
            <!-- 悬浮触发区域（隐形），保证展开后鼠标在缝隙间不失去焦点 -->
            <div class="reference-group-hover-trigger"></div>
            
            <!-- 已上传的图片卡片 -->
            <div class="preview-container-inline" v-for="(img, index) in form.reference_image" :key="index" :style="{ '--index': index, '--rotate': `${index === 0 ? 0 : [6, -4, 2, -8, 8, -6, 4, -2, 10, -10, 5, -5, 7, -7, 3, -3][(index - 1) % 16]}deg`, zIndex: form.reference_image.length - index + 1 }">
              <el-image 
                :src="img" 
                class="preview-img-inline" 
                :preview-src-list="windowWidth > 768 ? form.reference_image : []" 
                :initial-index="index"
                fit="cover"
                :preview-teleported="true"
                :hide-on-click-modal="true"
              />
              <div class="remove-btn-inline" @click.stop="clearReferenceImage(index)">
                <el-icon><Close /></el-icon>
              </div>
              <!-- 仅在第一张图片上显示的小型上传按钮（未展开时显示） -->
              <div class="collapsed-upload-badge" @click.stop="triggerUpload" v-if="index === 0 && form.reference_image.length < currentModelMaxRefImages">
                <svg width="1em" height="1em" viewBox="0 0 24 24" preserveAspectRatio="xMidYMid meet" fill="none" role="presentation" xmlns="http://www.w3.org/2000/svg"><g><path data-follow-fill="currentColor" d="M10.8 20a1.2 1.2 0 0 0 2.4 0v-6.8H20a1.2 1.2 0 1 0 0-2.4h-6.8V4a1.2 1.2 0 0 0-2.4 0v6.8H4a1.2 1.2 0 0 0 0 2.4h6.8V20Z" clip-rule="evenodd" fill-rule="evenodd" fill="currentColor"></path></g></svg>
              </div>
            </div>

            <!-- 折叠在卡片堆里的同尺寸大卡片续传按钮（仅悬浮时滑出） -->
            <div class="reference-upload-item-inline" @click.stop="triggerUpload" v-if="form.reference_image.length < currentModelMaxRefImages" :style="{ '--index': form.reference_image.length, '--rotate': `-12deg`, zIndex: 0 }">
              <div class="reference-upload-content" style="transform: rotate(12deg);">
                <svg width="1em" height="1em" viewBox="0 0 24 24" preserveAspectRatio="xMidYMid meet" fill="none" role="presentation" xmlns="http://www.w3.org/2000/svg" class="upload-icon-svg"><g><path data-follow-fill="currentColor" d="M10.8 20a1.2 1.2 0 0 0 2.4 0v-6.8H20a1.2 1.2 0 1 0 0-2.4h-6.8V4a1.2 1.2 0 0 0-2.4 0v6.8H4a1.2 1.2 0 0 0 0 2.4h6.8V20Z" clip-rule="evenodd" fill-rule="evenodd" fill="currentColor"></path></g></svg>
              </div>
            </div>
          </div>
        </div>
        
        <input type="file" ref="fileInputRef" accept="image/*" style="display: none" @change="handleFileUpload" multiple />
        <el-input
            ref="inputRef"
            v-model="form.prompt"
            type="textarea"
            :rows="1"
            :autosize="{ minRows: 1, maxRows: 6 }"
            placeholder="请输入你想生成的图片描述..."
            resize="none"
            class="chat-input"
            :class="{ 'has-references': form.reference_image && form.reference_image.length > 0, 'no-refs-allowed': currentModelMaxRefImages === 0 }"
            @keydown.enter.prevent="handleEnter"
            @focus="isInputFocused = true; isInputShrunk = false"
            @blur="isInputFocused = false"
            @click="isInputFocused = true; isInputShrunk = false"
          />
        <el-button
          v-if="windowWidth <= 768 && isMobileRefExpanded"
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
          :class="{ 'maintenance-btn': isCurrentResolutionMaintenance }"
          :disabled="!form.prompt.trim() || isCurrentResolutionMaintenance"
          @click="handleGenerate"
          :title="isCurrentResolutionMaintenance ? '该模型正在维护中，暂时无法生成' : '开始生成'"
        >
          <el-icon v-if="!isCurrentResolutionMaintenance"><Position /></el-icon>
          <span v-else style="font-size: 12px; transform: scale(0.9);">维护中</span>
        </el-button>
      </div>
    </div>
  </div>
</main>

<ExploreContent 
  v-if="route.path === '/explore'" 
  :form="form"
  :availableModels="availableModels"
  :currentModelInfo="currentModelInfo"
  :currentModelRatios="currentModelRatios"
  :currentModelResolutions="currentModelResolutions"
  :currentModelImageCounts="currentModelImageCounts"
  :currentModelMaxRefImages="currentModelMaxRefImages"
  :hasFreeResolution="hasFreeResolution"
  :getTierConfigForRes="getTierConfigForRes"
  :isCurrentResolutionMaintenance="isCurrentResolutionMaintenance"
  :windowWidth="windowWidth"
  :triggerUpload="triggerUpload"
  :clearReferenceImage="clearReferenceImage"
  :handleFileUpload="handleFileUpload"
  @toggle-menu="toggleMobileMenu" 
  @points-click="handleCommand('pointsRecord')" 
  @new-chat="handleNewChatFromExplore" 
  @generate="handleGenerate"
/>

<!-- 积分记录对话框 -->
<el-dialog v-model="showPointsRecordDialog" :width="400" destroy-on-close class="points-record-dialog" style="background: transparent; box-shadow: none;" :show-close="false" align-center append-to-body>
  <div class="points-container">
    <!-- 头部 -->
    <div class="points-header">
      <span class="points-title">积分中心</span>
      <div class="points-close-btn" @click="showPointsRecordDialog = false">
        <el-icon><Close /></el-icon>
      </div>
    </div>

    <!-- 余额卡片 -->
    <div class="points-balance-card">
      <div class="balance-info">
        <div style="display: flex; align-items: center; gap: 6px;">
          <span class="balance-label">当前剩余积分</span>
          <el-icon style="cursor: pointer; color: #9ca3af; font-size: 16px; margin-top: 1px;" @click="isPointsVisible = !isPointsVisible">
            <View v-if="isPointsVisible" />
            <Hide v-else />
          </el-icon>
        </div>
        <div class="balance-value">{{ isPointsVisible ? parseFloat(Number(userPoints || 0).toFixed(2)) : '****' }}</div>
      </div>
    </div>

    <!-- 明细列表 -->
    <div class="points-list-section" v-loading="pointsRecordLoading">
      <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
        <div class="list-section-title" style="margin-bottom: 0;">积分明细</div>
        <div class="balance-actions">
          <button class="balance-recharge-btn" @click="handleRecharge" v-if="siteSettings?.enable_recharge === 'true'">充值</button>
          <button class="balance-redeem-btn" @click="showRedeemDialog = true">卡密兑换</button>
        </div>
      </div>
      <div class="points-list-content">
        <template v-if="pointsRecordList && pointsRecordList.length > 0">
          <div class="points-list-item" v-for="(item, index) in pointsRecordList" :key="index">
            <div class="item-left">
              <div class="item-reason">{{ item.title || '官方发放' }}</div>
              <div class="item-reason-detail" style="font-size: 11px; color: #9ca3af; margin-top: 2px;">{{ item.reason }}</div>
              <div class="item-time">{{ new Date(item.created_at).toLocaleString('zh-CN', { hour12: false }).replace(/\//g, '-') }}</div>
            </div>
            <div class="item-right" :class="item.points_change > 0 ? 'is-positive' : 'is-negative'">
              {{ item.points_change > 0 ? '+' : '' }}{{ parseFloat(Number(item.points_change || 0).toFixed(2)) }}
            </div>
          </div>
        </template>
        <el-empty v-else description="暂无记录" :image-size="60" />
      </div>
      
      <div class="points-pagination-minimal" v-if="pointsRecordTotal > 0">
        <el-pagination
          v-model:current-page="pointsRecordPage"
          v-model:page-size="pointsRecordPageSize"
          :page-sizes="[10, 20, 50]"
          layout="prev, pager, next"
          :total="pointsRecordTotal"
          @size-change="handlePointsRecordSizeChange"
          @current-change="handlePointsRecordCurrentChange"
          small
        />
      </div>
    </div>
  </div>
</el-dialog>
<!-- 卡密兑换对话框 -->
<el-dialog v-model="showRedeemDialog" :width="420" destroy-on-close append-to-body class="premium-redeem-dialog" :show-close="false">
  <div class="premium-redeem-container">
    <div class="redeem-close-btn" @click="showRedeemDialog = false">
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path d="M1 1L13 13M1 13L13 1" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
      </svg>
    </div>
    
    <div class="redeem-header">
      <div class="redeem-icon-wrapper">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M20 12V22H4V12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M22 7H2V12H22V7Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M12 22V7" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M12 7H16.5C17.163 7 17.7989 6.73661 18.2678 6.26777C18.7366 5.79893 19 5.16304 19 4.5C19 3.83696 18.7366 3.20107 18.2678 2.73223C17.7989 2.26339 17.163 2 16.5 2C15.5 2 12 7 12 7Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M12 7H7.5C6.83696 7 6.20107 6.73661 5.73223 6.26777C5.26339 5.79893 5 5.16304 5 4.5C5 3.83696 5.26339 3.20107 5.73223 2.73223C6.20107 2.26339 6.83696 2 7.5 2C8.5 2 12 7 12 7Z" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      </div>
      <h3 class="redeem-title">兑换卡密</h3>
      <p class="redeem-desc">输入您的卡密，积分将直接充入您的账户</p>
    </div>
    
    <div class="redeem-body">
      <input v-model="cdkeyValue" class="custom-redeem-input" placeholder="请输入卡密兑换码" />
    </div>
    
    <div class="redeem-footer">
      <button class="redeem-action-btn cancel" @click="showRedeemDialog = false">取消</button>
      <button class="redeem-action-btn confirm" @click="handleRedeemCdkey" :disabled="redeemingCdkey">
        <span v-if="!redeemingCdkey">立即兑换</span>
        <span v-else>兑换中...</span>
      </button>
    </div>
  </div>
</el-dialog>

<!-- 移动到项目组弹窗 -->
<el-dialog v-model="moveToGroupDialogVisible" :width="420" destroy-on-close append-to-body align-center class="modern-group-dialog" :show-close="false">
  <div class="modern-dialog-header">
    <h3 class="modern-dialog-title">更换项目组</h3>
    <p class="modern-dialog-subtitle">将当前对话移动至指定的分组中进行管理</p>
  </div>
  
  <div class="modern-group-list-wrapper">
    <div class="modern-group-list" v-if="projectGroups.length > 0">
      <div 
        v-for="group in projectGroups" 
        :key="group.id" 
        class="modern-group-item"
        :class="{ 'is-active': targetProjectGroupId === group.id }"
        @click="targetProjectGroupId = group.id"
      >
        <div class="modern-group-item-left">
          <div class="modern-group-icon-wrapper">
            <el-icon class="modern-group-icon">
              <FolderOpened v-if="targetProjectGroupId === group.id" />
              <Folder v-else />
            </el-icon>
          </div>
          <span class="modern-group-name" :title="group.name">{{ group.name }}</span>
        </div>
        <div class="modern-check-wrapper" :class="{ 'is-checked': targetProjectGroupId === group.id }">
          <el-icon v-if="targetProjectGroupId === group.id"><Check /></el-icon>
        </div>
      </div>
    </div>
    <div v-else class="modern-group-empty">
      <el-empty description="暂无项目组，请先在侧边栏创建" :image-size="60"></el-empty>
    </div>
  </div>
  
  <div class="modern-dialog-footer">
    <button class="modern-btn modern-btn-cancel" @click="moveToGroupDialogVisible = false">取消</button>
    <button class="modern-btn modern-btn-confirm" @click="confirmMoveToGroup" :disabled="isMovingToGroup || targetProjectGroupId === null">
      {{ isMovingToGroup ? '移动中...' : '确认移动' }}
    </button>
  </div>
</el-dialog>

  </div>
</template>

<script setup lang="ts">
import { ref, reactive, nextTick, onMounted, computed, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Picture, MagicStick, Sunny, Star,
  Loading, Download, Position, Warning, Edit,
  Monitor, Brush, Close, MoreFilled, EditPen, Delete, Plus, Folder, FolderOpened, FolderRemove,
  Expand, Coin, Document, Cpu, Crop, CopyDocument, Clock, Operation, SwitchButton, List, ArrowRight, Setting, ArrowDown, RefreshRight, RefreshLeft, View, Hide, Check, Compass, CircleClose, UserFilled
} from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import { useAuthStore } from '@/store/auth'
import ExploreContent from '@/components/ExploreContent.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const contentAreaRef = ref<HTMLElement | null>(null)

// currentMainView not used anymore, removed

// Inspirations logic
const publicInspirations = ref<{id: number, content: string, image_url: string}[]>([])
const loadingInspirations = ref(false)

const fetchPublicInspirations = async () => {
  loadingInspirations.value = true
  try {
    const res: any = await request.get('/api/public/inspirations/list')
    publicInspirations.value = res.list || []
  } catch (error) {
    console.error('Failed to fetch inspirations', error)
  } finally {
    loadingInspirations.value = false
  }
}

// 滚动控制输入框缩放
const isInputShrunk = ref(false)
const lastScrollTop = ref(0)
const isAutoScrolling = ref(true)

const handleContentScroll = (e: Event) => {
  const target = e.target as HTMLElement
  const currentScrollTop = target.scrollTop
  
  // 忽略极小的滚动以防抖动
  if (Math.abs(currentScrollTop - lastScrollTop.value) < 10) return

  // 判断用户是否处于底部（误差 150px 以内认为在底部）
  isAutoScrolling.value = target.scrollHeight - currentScrollTop - target.clientHeight < 150

  // 只要发生有效滑动，就保持收缩状态并关闭移动端参考图展开
  isInputShrunk.value = true
  isMobileRefExpanded.value = false
  
  // 如果滚动到底部，恢复展开状态
  if (target.scrollHeight - currentScrollTop - target.clientHeight < 50) {
    isInputShrunk.value = false
  }

  lastScrollTop.value = currentScrollTop
}

const currentTime = ref(Date.now())
let timerInterval: number | null = null
import defaultLogo from '@/assets/logo/logo.png'

const cachedLogo = localStorage.getItem('site_logo')
const siteLogo = ref(cachedLogo || defaultLogo)
const isMobileMenuOpen = ref(false)
const isProjectGroupsExpanded = ref(true)
const isRecentTasksExpanded = ref(true)
const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

const currentSessionTitle = computed(() => {
  if (activeSessionId.value) {
    const session = sessionList.value.find(s => s.session_id === activeSessionId.value)
    if (session && session.title) {
      // 记录到 sessionStorage，防止刷新时出现标题空白
      sessionStorage.setItem('activeSessionTitle', session.title)
      return session.title
    }
    // 如果还没加载出来，尝试从缓存取
    const savedTitle = sessionStorage.getItem('activeSessionTitle')
    if (savedTitle) {
      return savedTitle
    }
  }
  return ''
})

const currentGroupName = computed(() => {
  if (activeSessionId.value) {
    const session = sessionList.value.find(s => s.session_id === activeSessionId.value)
    if (session && session.project_id && session.project_id !== 0) {
      const group = projectGroups.value.find(g => g.id === session.project_id)
      if (group) {
        sessionStorage.setItem('activeSessionGroupName', group.name)
        return group.name
      }
    }
    const savedGroupName = sessionStorage.getItem('activeSessionGroupName')
    if (savedGroupName && session && session.project_id !== 0) {
      return savedGroupName
    }
  }
  return ''
})

const username = computed(() => {
  const user = authStore.userInfo
  if (user) {
    return user.username || 'User'
  }
  return 'User'
})

const userEmail = computed(() => {
  const user = authStore.userInfo
  return user ? (user.email || '') : ''
})

const userPoints = computed(() => {
  const user = authStore.userInfo
  return user ? (user.points || 0) : 0
})
const isLoggedIn = computed(() => !!authStore.token)
const imageDomain = ref('')
const isInitialLoading = ref(true)

watch(() => authStore.token, (newToken, oldToken) => {
  if (newToken && !oldToken) {
    // User just logged in
    isInitialLoading.value = true // 显示骨架屏
    fetchUserInfo()
    fetchProjectGroups().then(() => {
      const fetchPromises = [fetchSessionsByGroup(0)]
      projectGroups.value.forEach(g => {
        fetchPromises.push(fetchSessionsByGroup(g.id))
      })
      
      Promise.all(fetchPromises).finally(() => {
        // 稍微延迟一下关闭骨架屏，防止闪烁太快
        setTimeout(() => {
          isInitialLoading.value = false
        }, 300)
      })
    })
  }
})
const windowWidth = ref(window.innerWidth)
const desktopPopoverRef = ref<any>(null)
const mobilePopoverRef = ref<any>(null)

let resizeTimer: ReturnType<typeof setTimeout> | null = null
const handleResize = () => {
  if (resizeTimer) clearTimeout(resizeTimer)
  resizeTimer = setTimeout(() => {
    windowWidth.value = window.innerWidth
    // 关闭可能处于打开状态的弹窗
    if (desktopPopoverRef.value) desktopPopoverRef.value.hide()
    if (mobilePopoverRef.value) mobilePopoverRef.value.hide()
  }, 150)
}

const fetchSettings = async () => {
  try {
    const res = await request.get('/api/public/settings')
    if (res.data && res.data.code === 200) {
      const settings = res.data.data || {}
      if (settings.image_domain !== undefined) {
        imageDomain.value = settings.image_domain.trim().replace(/\/$/, '')
        
        // 强制刷新历史记录和对话列表中的图片显示
        if (sessionList.value.length > 0) {
          sessionList.value = [...sessionList.value]
        }
        if (taskList.value.length > 0) {
          taskList.value = [...taskList.value]
        }
      }
      if (settings.site_name) {
        document.title = settings.site_name
      } else {
        document.title = 'Miren Al'
      }
      if (settings.site_logo !== undefined) {
        if (settings.site_logo) {
          siteLogo.value = settings.site_logo
          localStorage.setItem('site_logo', settings.site_logo)
          // Update favicon dynamically
          let link = document.querySelector("link[rel~='icon']") as HTMLLinkElement
          if (!link) {
            link = document.createElement('link')
            link.rel = 'icon'
            document.head.appendChild(link)
          }
          link.href = settings.site_logo
        } else {
          siteLogo.value = defaultLogo
          localStorage.removeItem('site_logo')
          let link = document.querySelector("link[rel~='icon']") as HTMLLinkElement
          if (link) {
            link.href = '/logo.png'
          }
        }
      }
    }
  } catch (error) {
    console.warn('Failed to fetch public settings', error)
    document.title = 'Miren Al'
    ElMessage.error('无法连接到服务器，请检查网络连接或稍后重试。')
  }
}

const taskList = ref<any[]>([])
const deletedTaskIds = ref<Set<number>>(new Set())
const sessionList = ref<any[]>([])
const projectGroups = ref<any[]>([])
const projectGroupsPagination = ref({ page: 1, hasMore: false, loadingMore: false })
const collapsedGroups = ref<Set<number>>(new Set())

const toggleGroup = (groupId: number) => {
  if (collapsedGroups.value.has(groupId)) {
    collapsedGroups.value.delete(groupId)
  } else {
    collapsedGroups.value.add(groupId)
  }
}

const openDropdownId = ref<string | number | null>(null)
const handleDropdownVisible = (visible: boolean, id: string | number) => {
  if (visible) {
    openDropdownId.value = id
  } else if (openDropdownId.value === id) {
    openDropdownId.value = null
  }
}

const editingGroupId = ref<number | null>(null)
const editGroupName = ref('')
const editGroupInputRef = ref<any>(null)

const sessionsByGroup = computed(() => {
  const map: Record<number, any[]> = {}
  map[0] = []
  if (projectGroups.value) {
    projectGroups.value.forEach(g => {
      map[g.id] = []
    })
  }
  if (sessionList.value) {
    sessionList.value.forEach(s => {
      const gid = s.project_id || 0
      if (!map[gid]) {
        map[gid] = []
      }
      map[gid].push(s)
    })
  }
  return map
})

const getSessionsByGroupId = (groupId: number) => {
  return sessionsByGroup.value[groupId] || []
}

const shouldShowChatTime = (task: any, index: number) => {
  if (index === 0) return true
  const prevTask = chatList.value[index - 1]
  if (!prevTask || !prevTask.created_at || !task.created_at) return true
  
  const currTime = new Date(task.created_at).getTime()
  const prevTime = new Date(prevTask.created_at).getTime()
  
  // 两个消息相差超过5分钟 (300000 毫秒) 则显示时间
  return (currTime - prevTime) > 300000
}

const formatChatTime = (timeStr: string) => {
  if (!timeStr) return ''
  // 转换 UTC 时间到本地时间 (东八区)
  // 后端返回的是 2026-09-06T15:21:28Z，new Date() 能够自动根据 Z 解析为本地时间
  const date = new Date(timeStr)
  if (isNaN(date.getTime())) return timeStr
  
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const yesterday = new Date(today.getTime() - 86400000)
  
  const targetDate = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  
  const pad = (n: number) => (n < 10 ? '0' + n : n)
  const timePart = `${pad(date.getHours())}:${pad(date.getMinutes())}`
  
  if (targetDate.getTime() === today.getTime()) {
    return `今天 ${timePart}`
  } else if (targetDate.getTime() === yesterday.getTime()) {
    return `昨天 ${timePart}`
  } else if (date.getFullYear() === now.getFullYear()) {
    return `${date.getMonth() + 1}月${date.getDate()}日 ${timePart}`
  } else {
    return `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日 ${timePart}`
  }
}

const formatTime = (timeStr: string) => {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  if (isNaN(date.getTime())) return timeStr
  
  const pad = (n: number) => (n < 10 ? '0' + n : n)
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

const activePopoverId = ref<string | number | null>(null)
const globalGeneratingGroupIds = ref<number[]>([])

const copyPrompt = async (text: string) => {
  if (!text) return
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('提示词已复制')
  } catch (err) {
    ElMessage.error('复制失败，请手动复制')
  }
}

const fetchProjectGroups = async (loadMore = false, silent = false) => {
  if (loadMore && !projectGroupsPagination.value.hasMore) return
  if (loadMore) projectGroupsPagination.value.loadingMore = true

  try {
    const wasEmpty = projectGroups.value.length === 0
    const currentPage = projectGroupsPagination.value.page
    const pageToFetch = loadMore ? currentPage + 1 : 1
    
    let pageSize = 5
    if (!loadMore && currentPage > 1) {
      pageSize = currentPage * 5
    }

    const res = await request.get(`/api/user/projects/list?page=${pageToFetch}&page_size=${pageSize}`)
    if (res.data && res.data.code === 200 && res.data.data) {
      const data = res.data.data
      const list = data.list || []
      globalGeneratingGroupIds.value = data.generating_group_ids || []
      
      // 预先发起并行请求，让对话数据尽早开始拉取
      let fetchPromises: Promise<any>[] = []
      if (!silent) {
        fetchPromises = list.map((g: any) => {
          const hasSessions = sessionList.value.some(s => s.project_id === g.id)
          if (!hasSessions) {
            return fetchSessionsByGroup(g.id, false, false)
          }
          return Promise.resolve()
        })
      }
      
      if (loadMore) {
        // 使用服务器返回的完整列表顺序进行替换，并保持分页状态
        const newItems = list.filter((g: any) => !projectGroups.value.some(exist => exist.id === g.id))
        projectGroups.value.push(...newItems)
        newItems.forEach((g: any) => {
          if (!collapsedGroups.value.has(g.id)) {
            collapsedGroups.value.add(g.id)
          }
        })
        projectGroupsPagination.value.page = pageToFetch
        
        // 只有当这是点击加载更多，并且返回的数据确实表明没有下一页时，才标记没有更多了。
        // 由于 loadMore 时，请求的只是“下一页”的数据，所以 data.total 应该是整体的条数，
        // 而当前已加载的总条数就是 pageToFetch * pageSize (这里是 5)。
        if (list.length < pageSize || data.total <= pageToFetch * pageSize) {
          projectGroupsPagination.value.hasMore = false
        }
        // 如果这是初始加载或刷新，而不是点击“加载更多”
      } else {
        // 查找当前选中的会话
        const currentActiveSession = sessionList.value.find(s => s.session_id === activeSessionId.value)
        // 如果当前选中的会话在一个项目组里，我们要确保这个项目组不会被刷掉
        let activeGroupId: number | null = null
        if (currentActiveSession && currentActiveSession.project_id && currentActiveSession.project_id !== 0) {
          activeGroupId = currentActiveSession.project_id
        }
        
        const existingActiveGroup = activeGroupId ? projectGroups.value.find(g => g.id === activeGroupId) : null
        
        if (!silent) {
          projectGroups.value = list
          if (existingActiveGroup && !list.some((g: any) => g.id === existingActiveGroup.id)) {
            projectGroups.value.push(existingActiveGroup)
          }
        } else {
          // 静默刷新，直接使用服务器返回的列表和顺序，并保留其他特殊组
          // 但是要注意：静默刷新拉取的是 1 到 N 页的数据，我们要保留在这 N 页之外（可能更深处）已经被我们加载并激活的组
          const updatedCurrentGroups = list.map((serverGroup: any) => {
            const existing = projectGroups.value.find(g => g.id === serverGroup.id)
            return existing ? { ...existing, ...serverGroup } : serverGroup
          })
          
          // 确保 activeGroupId 被保留，即使它不在当前拉取的列表中
          if (existingActiveGroup && !list.some((g: any) => g.id === existingActiveGroup.id)) {
            updatedCurrentGroups.push(existingActiveGroup)
          }
          
          projectGroups.value = updatedCurrentGroups
        }
        
        if (!silent && wasEmpty) {
          const allGroupIds = list.map((g: any) => g.id)
          collapsedGroups.value = new Set(allGroupIds)
        }
        
        // 等待所有对话数据拉取完成
        if (fetchPromises.length > 0) {
          await Promise.allSettled(fetchPromises)
        }
        
        // 静默刷新或首次加载时，当前请求的数据量是 pageSize（已乘以 page）
        // 如果服务器返回的总数小于等于我们请求的容量，那说明没有更多了。
        if (list.length < pageSize || data.total <= pageSize) {
          projectGroupsPagination.value.hasMore = false
        } else {
          projectGroupsPagination.value.hasMore = true
        }
      }
    }
  } catch (error) {
    console.error('Failed to fetch project groups', error)
  } finally {
    if (loadMore) projectGroupsPagination.value.loadingMore = false
  }
}

const createProjectGroup = async () => {
  if (!isLoggedIn.value) {
    authStore.openLogin()
    return
  }
  ElMessageBox.prompt('请输入项目组名称', '新建项目组', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    inputPattern: /\S+/,
    inputErrorMessage: '名称不能为空',
    inputValidator: (value) => {
      if (value && value.trim().length > 10) {
        return '项目组名称不能超过10个字符'
      }
      return true
    }
  }).then(async ({ value }) => {
    try {
      const res = await request.post('/api/user/projects/create', { name: value.trim() })
      if (res.data && res.data.code === 200) {
        ElMessage.success('创建成功')
        fetchProjectGroups()
      } else {
        ElMessage.error(res.data.message || '创建失败')
      }
    } catch (error) {
      ElMessage.error('创建失败')
    }
  }).catch(() => {})
}



const handleGroupCommand = (command: string, group: any) => {
  if (command === 'rename') {
    editingGroupId.value = group.id
    editGroupName.value = group.name
    nextTick(() => {
      if (editGroupInputRef.value && editGroupInputRef.value.length > 0) {
        editGroupInputRef.value[0].focus()
      } else if (editGroupInputRef.value) {
        editGroupInputRef.value.focus()
      }
    })
  } else if (command === 'delete') {
    ElMessageBox.confirm('确定要删除该项目组吗？组内的任务不会被删除，将移回"最近任务"。', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }).then(async () => {
      try {
        const res = await request.post('/api/user/projects/delete', { id: group.id })
        if (res.data && res.data.code === 200) {
            ElMessage.success('删除成功')
            // Optimistic update for better performance
            projectGroups.value = projectGroups.value.filter(g => g.id !== group.id)
            sessionList.value = sessionList.value.map(s => {
              if (s.project_id === group.id) {
                return { ...s, project_id: 0 }
              }
              return s
            })
            // Fetch recent tasks to ensure pagination is consistent
            fetchSessionsByGroup(0, false, true)
          } else {
          ElMessage.error(res.data.message || '删除失败')
        }
      } catch (error) {
        ElMessage.error('删除失败')
      }
    }).catch(() => {})
  }
}

const handleGroupNameInput = (value: string) => {
  if (value.length > 10) {
    ElMessage.warning('项目组名称不能超过10个字符')
    editGroupName.value = value.slice(0, 10)
  }
}

const saveGroupName = async (group: any) => {
  if (!editingGroupId.value) return
  
  const newName = editGroupName.value.trim()
  editingGroupId.value = null
  
  if (newName && newName !== group.name) {
    if (newName.length > 10) {
      ElMessage.warning('项目组名称不能超过10个字符')
      return
    }
    try {
      const res = await request.post('/api/user/projects/update', { 
        id: group.id,
        name: newName
      })
      if (res.data && res.data.code === 200) {
        group.name = newName
        ElMessage.success('重命名成功')
      } else {
        ElMessage.error(res.data.message || '重命名失败')
      }
    } catch (error) {
      ElMessage.error('重命名失败')
    }
  }
}

const handleDragStartSession = (event: DragEvent, session: any) => {
  if (event.dataTransfer) {
    event.dataTransfer.setData('text/plain', session.session_id)
    event.dataTransfer.effectAllowed = 'move'
  }
}

// 侧边栏拖拽自动滚动逻辑
const sidebarContentRef = ref<HTMLElement | null>(null)
let dragScrollFrame: number | null = null
let currentDragY = 0

const handleSidebarDragOver = (event: DragEvent) => {
  event.preventDefault()
  currentDragY = event.clientY
  
  if (!dragScrollFrame) {
    startDragScroll()
  }
}

const handleSidebarDragLeave = (event: DragEvent) => {
  // 如果鼠标离开了整个侧边栏区域，停止滚动
  const container = sidebarContentRef.value
  if (!container) return
  const rect = container.getBoundingClientRect()
  if (
    event.clientX < rect.left || 
    event.clientX > rect.right || 
    event.clientY < rect.top || 
    event.clientY > rect.bottom
  ) {
    stopDragScroll()
  }
}

const startDragScroll = () => {
  if (!sidebarContentRef.value) return
  
  const container = sidebarContentRef.value
  const rect = container.getBoundingClientRect()
  
  const topDistance = currentDragY - rect.top
  const bottomDistance = rect.bottom - currentDragY
  const scrollThreshold = 100 // 扩大边缘触发区域
  const maxScrollSpeed = 20 // 增加滚动速度
  
  let scrolled = false
  
  if (topDistance < scrollThreshold && topDistance > -20) { // 稍微允许超出一点
    const speed = maxScrollSpeed * (1 - Math.max(0, topDistance) / scrollThreshold)
    container.scrollTop -= Math.max(speed, 3)
    scrolled = true
  } else if (bottomDistance < scrollThreshold && bottomDistance > -20) {
    const speed = maxScrollSpeed * (1 - Math.max(0, bottomDistance) / scrollThreshold)
    container.scrollTop += Math.max(speed, 3)
    scrolled = true
  }
  
  if (scrolled) {
    dragScrollFrame = requestAnimationFrame(startDragScroll)
  } else {
    dragScrollFrame = null
  }
}

const stopDragScroll = () => {
  if (dragScrollFrame) {
    cancelAnimationFrame(dragScrollFrame)
    dragScrollFrame = null
  }
}

const isMutatingSession = ref(false)

const handleDropOnSession = async (event: DragEvent, targetSession: any, targetIndex: number, projectId: number) => {
  stopDragScroll()
  event.preventDefault()
  if (!event.dataTransfer) return
  
  const sourceSessionId = event.dataTransfer.getData('text/plain')
  if (!sourceSessionId || sourceSessionId === targetSession.session_id) return
  
  const sourceSession = sessionList.value.find(s => s.session_id === sourceSessionId)
  if (!sourceSession) return

  isMutatingSession.value = true

  // 获取目标分组内的所有会话
  const groupSessions = getSessionsByGroupId(projectId)
  const sourceIndex = groupSessions.findIndex(s => s.session_id === sourceSessionId)

  // 如果是在同一个组内拖拽，或者从其他组跨组拖入
  let newGroupSessions = [...groupSessions]
  let isCrossGroup = false
  const originalProjectId = sourceSession.project_id || 0
  
  if ((sourceSession.project_id || 0) === projectId) {
    // 同组内排序
    newGroupSessions.splice(sourceIndex, 1)
    newGroupSessions.splice(targetIndex, 0, sourceSession)
  } else {
    // 跨组移动并排序
    isCrossGroup = true
    sourceSession.project_id = projectId
    newGroupSessions.splice(targetIndex, 0, sourceSession)
    if (projectId !== 0 && collapsedGroups.value.has(projectId)) {
      collapsedGroups.value.delete(projectId)
    }
  }

  // 立即完成前端状态更新（乐观更新），不等待网络请求
  const otherSessions = sessionList.value.filter(s => (s.project_id || 0) !== projectId)
  
  // 对于跨组移动，我们必须从其原始的组列表中将它移除。
  // 因为 groupSessions 只是目标组的副本，而 otherSessions 会包含所有其他组（包括原组）。
  // 当我们执行 sourceSession.project_id = projectId 之后，
  // otherSessions 的 filter 条件 (s.project_id || 0) !== projectId 会将它排除掉。
  // 这确保了它不会在 otherSessions 中重复出现，并且由于我们在 newGroupSessions 中添加了它，
  // 所以整个 sessionList 仍然是完整的，只是它在原组的“残留”被干净地清除了。
  sessionList.value = [...otherSessions, ...newGroupSessions]
  
  // 如果是跨组移动，发送移动请求
  if (isCrossGroup) {
    try {
      const resMove = await request.post('/api/user/sessions/move', { 
        session_id: sourceSessionId,
        project_id: projectId
      })
      if (resMove.data && resMove.data.code !== 200) {
        throw new Error(resMove.data.message || '移动失败')
      }
    } catch (e) {
      console.error('Failed to move session', e)
      // 如果移动失败，进行回滚
      sourceSession.project_id = originalProjectId
      sessionList.value = [...sessionList.value]
      ElMessage.error('移动并排序失败')
      isMutatingSession.value = false
      return
    }
  }

  // 将新的顺序发送给后端
  const orderedSessionIds = newGroupSessions.map(s => s.session_id)
  try {
    const res = await request.post('/api/user/sessions/reorder', {
      session_ids: orderedSessionIds
    })
    
    if (res.data && res.data.code === 200) {
      // 成功时不显示提示，保持乐观更新的无缝体验
    } else {
      // 在这里添加强制响应式更新
      sessionList.value = [...sessionList.value]
      ElMessage.error('排序失败')
    }
  } catch (error) {
    console.error('Failed to reorder sessions', error)
    // 失败时强制更新 UI，防止卡死
    sessionList.value = [...sessionList.value]
  } finally {
    isMutatingSession.value = false
  }
}

const handleDropSession = async (event: DragEvent, projectId: number) => {
  stopDragScroll()
  event.preventDefault()
  if (!event.dataTransfer) return
  
  const sessionId = event.dataTransfer.getData('text/plain')
  if (!sessionId) return
  
  const session = sessionList.value.find(s => s.session_id === sessionId)
  if (!session || (session.project_id || 0) === projectId) return
  
  isMutatingSession.value = true

  // Optimistic update
  const originalProjectId = session.project_id || 0
  session.project_id = projectId
  
  // Force update sessionList reference to trigger reactivity in computed properties
  sessionList.value = [...sessionList.value]
  
  if (projectId !== 0 && collapsedGroups.value.has(projectId)) {
    collapsedGroups.value.delete(projectId)
  }
  
  try {
    const res = await request.post('/api/user/sessions/move', { 
      session_id: sessionId,
      project_id: projectId
    })
    
    if (res.data && res.data.code === 200) {
      // 成功时不显示提示，保持乐观更新的无缝体验
    } else {
      // Revert on failure
      session.project_id = originalProjectId
      sessionList.value = [...sessionList.value]
      ElMessage.error(res.data.message || '移动失败')
    }
  } catch (error) {
    // Revert on failure
    session.project_id = originalProjectId
    sessionList.value = [...sessionList.value]
    ElMessage.error('移动失败')
  } finally {
    isMutatingSession.value = false
  }
}
const availableModels = ref<any[]>([])

const getModelName = (seriesId: string) => {
  if (!seriesId || seriesId === 'default') return '默认模型'
  const model = availableModels.value.find(m => m.series_id === seriesId)
  return model ? model.name : seriesId
}

const getModelCredits = (seriesId: string, resolution?: string) => {
  if (!seriesId || seriesId === 'default') return 1
  const model = availableModels.value.find(m => m.series_id === seriesId)
  if (!model) return 1
  
  if (resolution && model.resolution_configs) {
    try {
      const configs = JSON.parse(model.resolution_configs)
      if (configs && configs[resolution] && configs[resolution].credits_per_image !== undefined && configs[resolution].credits_per_image !== '') {
        return Number(configs[resolution].credits_per_image)
      }
    } catch (e) {}
  }
  return model.credits_per_image !== undefined ? Number(model.credits_per_image) : 1
}

const calculateRefundedPoints = (task: any) => {
  if (task.status !== 2 && task.status !== 3) return 0
  const logContent = task.log_content || ''
  const partialMatch = logContent.match(/Partially generated \d+ images\. (\d+) failed\./)
  if (partialMatch && partialMatch[1]) {
    const failedCount = parseInt(partialMatch[1], 10)
    const creditsPerImage = getModelCredits(task.series_id, task.resolution)
    return parseFloat((Number(creditsPerImage) * failedCount).toFixed(2))
  }
  if (task.status === 2) {
    return parseFloat((getModelCredits(task.series_id, task.resolution) * (task.num_images || 1)).toFixed(2))
  }
  return 0
}

const tempTasksMap = ref<Record<string, any[]>>({})

// 从 sessionStorage 恢复临时任务（防止刷新后生成中的任务丢失）
try {
  const savedTemp = sessionStorage.getItem('tempTasksMap')
  if (savedTemp) {
    tempTasksMap.value = JSON.parse(savedTemp)
  }
} catch (e) {}

watch(tempTasksMap, (newVal) => {
  sessionStorage.setItem('tempTasksMap', JSON.stringify(newVal))
}, { deep: true })

const activeSessionId = ref<string | null>(null)


const editingSessionId = ref<string | null>(null)
const editSessionName = ref('')
const editInputRef = ref<any>(null)

const moveToGroupDialogVisible = ref(false)
const targetProjectGroupId = ref<number | null>(null)
const targetSessionForMove = ref<any>(null)
const isMovingToGroup = ref(false)

const confirmMoveToGroup = async () => {
  if (!targetSessionForMove.value || targetProjectGroupId.value === null) return
  isMovingToGroup.value = true
  const sourceSession = targetSessionForMove.value
  const projectId = targetProjectGroupId.value
  const originalProjectId = sourceSession.project_id || 0
  
  // 乐观更新（立即修改状态，消除延迟感）
  // 必须确保是从 sessionList 中找到并修改，因为 sourceSession 可能是个旧引用
  const sessionIndex = sessionList.value.findIndex(s => s.session_id === sourceSession.session_id)
  if (sessionIndex !== -1) {
    sessionList.value[sessionIndex] = { ...sessionList.value[sessionIndex], project_id: projectId }
    sessionList.value = [...sessionList.value]
  } else {
    sourceSession.project_id = projectId
    sessionList.value = [...sessionList.value]
  }

  if (projectId !== 0 && collapsedGroups.value.has(projectId)) {
    collapsedGroups.value.delete(projectId)
  }
  moveToGroupDialogVisible.value = false
  
  try {
    const resMove = await request.post('/api/user/sessions/move', { 
      session_id: sourceSession.session_id,
      project_id: projectId
    })
    if (resMove.data && resMove.data.code === 200) {
      ElMessage.success('移动成功')
    } else {
      // 失败则回滚状态
      if (sessionIndex !== -1) {
        sessionList.value[sessionIndex] = { ...sessionList.value[sessionIndex], project_id: originalProjectId }
      } else {
        sourceSession.project_id = originalProjectId
      }
      sessionList.value = [...sessionList.value]
      ElMessage.error(resMove.data?.message || '移动失败')
    }
  } catch (e) {
    // 失败则回滚状态
    if (sessionIndex !== -1) {
      sessionList.value[sessionIndex] = { ...sessionList.value[sessionIndex], project_id: originalProjectId }
    } else {
      sourceSession.project_id = originalProjectId
    }
    sessionList.value = [...sessionList.value]
    ElMessage.error('移动失败')
  } finally {
    isMovingToGroup.value = false
  }
}

const handleSessionCommand = (command: string, session: any) => {
  if (command === 'rename') {
    editingSessionId.value = session.session_id
    editSessionName.value = session.title || '图片生成任务'
    nextTick(() => {
      if (editInputRef.value && editInputRef.value.length > 0) {
        editInputRef.value[0].focus()
      } else if (editInputRef.value) {
        editInputRef.value.focus()
      }
    })
  } else if (command === 'moveToGroup') {
    targetSessionForMove.value = session
    targetProjectGroupId.value = session.project_id || null
    moveToGroupDialogVisible.value = true
  } else if (command === 'removeFromGroup') {
    targetSessionForMove.value = session
    targetProjectGroupId.value = 0
    confirmMoveToGroup()
  } else if (command === 'delete') {
    ElMessageBox.confirm('确定要删除该任务吗？此操作将同时删除该任务下的所有图片和记录，且不可恢复。', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }).then(async () => {
      // Optimistic update
      const backupList = [...sessionList.value]
      sessionList.value = sessionList.value.filter(s => s.session_id !== session.session_id)
      if (activeSessionId.value === session.session_id) {
        resetTask()
      }
      try {
        const res = await request.post('/api/user/sessions/delete', { session_id: session.session_id })
        if (res.data && res.data.code === 200) {
          ElMessage.success('删除成功')
        } else {
          // Revert on failure
          sessionList.value = backupList
          ElMessage.error(res.data.message || '删除失败')
        }
      } catch (error) {
        // Revert on failure
        sessionList.value = backupList
        ElMessage.error('删除失败')
      }
    }).catch(() => {})
  }
}

const saveSessionName = async (session: any) => {
  if (!editingSessionId.value) return
  
  const newName = editSessionName.value.trim()
  editingSessionId.value = null
  
  if (newName && newName !== session.title) {
    try {
      const res = await request.post('/api/user/sessions/rename', { 
        session_id: session.session_id,
        name: newName
      })
      if (res.data && res.data.code === 200) {
        session.title = newName
        ElMessage.success('重命名成功')
      } else {
        ElMessage.error(res.data.message || '重命名失败')
      }
    } catch (error) {
      ElMessage.error('重命名失败')
    }
  }
}

const isNewTaskMode = ref(true)

const currentModelInfo = computed(() => {
  if (form.series_id === 'default') {
    return availableModels.value.length > 0 ? availableModels.value[0] : null
  }
  return availableModels.value.find(m => m.series_id === form.series_id) || null
})

const getTierConfig = (info: any, field: string, fallback: any) => {
  if (!info) return fallback
  if (info.resolution_configs) {
    try {
      const configs = JSON.parse(info.resolution_configs)
      if (configs && configs[form.resolution] && configs[form.resolution][field] !== undefined && configs[form.resolution][field] !== '') {
        return configs[form.resolution][field]
      }
    } catch (e) {}
  }
  if (info[field] !== undefined && info[field] !== '') return info[field]
  return fallback
}

const getTierConfigForRes = (info: any, res: string, field: string, fallback: any) => {
  if (!info) return fallback
  if (info.resolution_configs) {
    try {
      const configs = JSON.parse(info.resolution_configs)
      if (configs && configs[res] && configs[res][field] !== undefined && configs[res][field] !== '') {
        return configs[res][field]
      }
    } catch (e) {}
  }
  if (info[field] !== undefined && info[field] !== '') return info[field]
  return fallback
}

const hasFreeResolution = (modelInfo: any) => {
  if (!modelInfo) return false
  if (modelInfo.resolution_tiers) {
    let tiers = String(modelInfo.resolution_tiers).split(',').filter((r: string) => r.trim())
    
    if (modelInfo.resolution_configs) {
      try {
        const configs = JSON.parse(modelInfo.resolution_configs)
        tiers = tiers.filter(tier => !(configs[tier] && configs[tier].enabled === false))
      } catch (e) {}
    }
    
    for (const res of tiers) {
      if (Number(getTierConfigForRes(modelInfo, res, 'credits_per_image', 1)) === 0) {
        return true
      }
    }
  } else {
     if (Number(modelInfo.credits_per_image) === 0) return true
  }
  return false
}

const currentModelRatios = computed(() => {
  const ratios = getTierConfig(currentModelInfo.value, 'aspect_ratios', null)
  if (ratios) {
    return String(ratios).split(',').filter((r: string) => r.trim())
  }
  return ['auto', '1:1', '4:3', '3:4', '3:2', '2:3', '16:9', '9:16', '21:9', '9:21'] // default fallback
})

const currentModelImageCounts = computed(() => {
  const counts = getTierConfig(currentModelInfo.value, 'image_counts', null)
  if (counts) {
    const strCounts = String(counts).trim()
    if (/^\d+$/.test(strCounts)) {
      const max = parseInt(strCounts, 10)
      if (max > 0) {
        return Array.from({ length: max }, (_, i) => String(i + 1))
      }
    }
    return strCounts.split(',').filter((c: string) => c.trim())
  }
  return ['1', '2', '3', '4'] // default fallback
})

const currentModelResolutions = computed(() => {
  if (currentModelInfo.value && currentModelInfo.value.resolution_tiers) {
    let tiers = String(currentModelInfo.value.resolution_tiers).split(',').filter((r: string) => r.trim())
    
    // Filter out disabled tiers based on resolution_configs
    if (currentModelInfo.value.resolution_configs) {
      try {
        const configs = JSON.parse(currentModelInfo.value.resolution_configs)
        tiers = tiers.filter(tier => !(configs[tier] && configs[tier].enabled === false))
      } catch (e) {}
    }
    
    return tiers
  }
  return ['1K', '2K', '4K'] // default fallback
})

const isCurrentResolutionMaintenance = computed(() => {
  return getTierConfigForRes(currentModelInfo.value, form.resolution, 'enabled', true) === 'maintenance'
})

const currentModelMaxRefImages = computed(() => {
  return Number(getTierConfig(currentModelInfo.value, 'max_reference_images', 3))
})

const generatingSessions = computed(() => {
  const sessions: string[] = []
  for (const [sessionId, tasks] of Object.entries(tempTasksMap.value)) {
    if (tasks.some(t => t.status === 0)) {
      sessions.push(sessionId)
    }
  }
  return sessions
})

const generatingCountByGroup = computed(() => {
  const map: Record<number, number> = {}
  map[0] = 0
  if (projectGroups.value) {
    projectGroups.value.forEach(g => {
      map[g.id] = 0
    })
  }
  
  Object.keys(sessionsByGroup.value).forEach(groupIdStr => {
    const groupId = Number(groupIdStr)
    const sessions = sessionsByGroup.value[groupId]
    let count = 0
    sessions.forEach(session => {
      let backendGenerating = session.generating_count || 0
      const tempTasks = tempTasksMap.value[session.session_id] || []
      let localGenerating = 0
      tempTasks.forEach(t => {
        if (t.status === 0) {
          localGenerating++
          if (t.id && !String(t.id).startsWith('temp-') && backendGenerating > 0) {
            backendGenerating--
          }
        }
      })
      count += backendGenerating + localGenerating
    })
    map[groupId] = count
  })
  return map
})

const getGeneratingCountInGroup = (groupId: number) => {
  return generatingCountByGroup.value[groupId] || 0
}

const generatingProjectGroupsCount = computed(() => {
  const allActiveGroupIds = new Set(globalGeneratingGroupIds.value)
  
  // 1. Add groups from optimistic state
  Object.keys(tempTasksMap.value).forEach(sessionId => {
    const tempTasks = tempTasksMap.value[sessionId]
    if (tempTasks.some(t => t.status === 0)) {
      const session = sessionList.value.find(s => s.session_id === sessionId)
      if (session && session.project_id && session.project_id > 0) {
        allActiveGroupIds.add(session.project_id)
      }
    }
  })

  // 2. Remove groups that we definitively know are NOT generating
  // For any group that we have loaded, we calculate its accurate local generating count.
  // If it's 0, we remove it from the Set to reflect immediate completion.
  projectGroups.value.forEach(g => {
    if (getGeneratingCountInGroup(g.id) === 0) {
      allActiveGroupIds.delete(g.id)
    } else {
      allActiveGroupIds.add(g.id)
    }
  })
  
  return allActiveGroupIds.size
})

const chatList = computed(() => {
  const currentTempTasks = activeSessionId.value ? (tempTasksMap.value[activeSessionId.value] || []) : []
  const taskIds = new Set(taskList.value.map(t => t.id))
  const uniqueTempTasks = currentTempTasks.filter(t => !taskIds.has(t.id))
  return [...taskList.value, ...uniqueTempTasks].filter(t => !deletedTaskIds.value.has(t.id))
})

const form = reactive({
  prompt: '',
  series_id: 'default',
  resolution: '1K',
  aspect_ratio: 'auto',
  num_images: 1,
  reference_image: [] as string[]
})

// 跟踪输入框是否获得焦点
const isInputFocused = ref(false)



// 跟踪移动端展开状态
const isMobileRefExpanded = ref(false)

// Click outside to remove focus state
const handleClickOutside = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest('.input-box')) {
    isInputFocused.value = false
  }
  // 点击外部时收起移动端的图片堆叠
  if (!target.closest('.reference-image-preview-inline')) {
    isMobileRefExpanded.value = false
  }
}

const pendingReferenceFiles = ref<File[]>([])

const clearReferenceImage = (index?: number) => {
  if (index !== undefined) {
    form.reference_image.splice(index, 1)
    pendingReferenceFiles.value.splice(index, 1)
  } else {
    form.reference_image = []
    pendingReferenceFiles.value = []
  }
}

const fileInputRef = ref<HTMLInputElement | null>(null)

const triggerUpload = () => {
  if (form.reference_image.length >= currentModelMaxRefImages.value) {
    ElMessage.warning(`最多只能上传 ${currentModelMaxRefImages.value} 张参考图`)
    return
  }
  fileInputRef.value?.click()
}

const handleFileUpload = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = target.files
  if (!files || files.length === 0) return

  let addedCount = 0;
  for (let i = 0; i < files.length; i++) {
    const file = files[i];
    
    if (file.size > 5 * 1024 * 1024) {
      ElMessage.warning(`图片 ${file.name} 大小不能超过5MB，已跳过`)
      continue;
    }
    
    if (form.reference_image.length >= currentModelMaxRefImages.value) {
      ElMessage.warning(`最多只能上传 ${currentModelMaxRefImages.value} 张参考图，已截断多余图片`)
      break;
    }

    pendingReferenceFiles.value.push(file)
    form.reference_image.push(URL.createObjectURL(file))
    addedCount++;
  }
  
  // reset input
  target.value = ''
  
  // 移动端：上传完成后自动展开，方便用户看到刚上传的图片并进行管理
  if (windowWidth.value <= 768) {
    isMobileRefExpanded.value = true
  }
}

watch(() => form.series_id, () => {
  // Update form values to be valid for the newly selected model
  // 如果切换模型后，原本选中的分辨率在新模型中不存在，则强制选中第一个分辨率
  if (currentModelResolutions.value.length > 0 && !currentModelResolutions.value.includes(form.resolution)) {
    form.resolution = currentModelResolutions.value[0]
  }
  
  if (currentModelRatios.value.length > 0 && !currentModelRatios.value.includes(form.aspect_ratio)) {
    form.aspect_ratio = currentModelRatios.value[0]
  }
  if (currentModelImageCounts.value.length > 0 && !currentModelImageCounts.value.includes(String(form.num_images))) {
    form.num_images = Number(currentModelImageCounts.value[0])
  }
  if (currentModelMaxRefImages.value === 0) {
    clearReferenceImage()
  } else if (form.reference_image.length > currentModelMaxRefImages.value) {
    form.reference_image.splice(currentModelMaxRefImages.value)
    pendingReferenceFiles.value.splice(currentModelMaxRefImages.value)
  }
})

watch(() => form.resolution, () => {
  if (currentModelRatios.value.length > 0 && !currentModelRatios.value.includes(form.aspect_ratio)) {
    form.aspect_ratio = currentModelRatios.value[0]
  }
  if (currentModelImageCounts.value.length > 0 && !currentModelImageCounts.value.includes(String(form.num_images))) {
    form.num_images = Number(currentModelImageCounts.value[0])
  }
  if (currentModelMaxRefImages.value === 0) {
    clearReferenceImage()
  } else if (form.reference_image.length > currentModelMaxRefImages.value) {
    form.reference_image.splice(currentModelMaxRefImages.value)
    pendingReferenceFiles.value.splice(currentModelMaxRefImages.value)
  }
})

let pollInterval: number | null = null

const fetchUserInfo = async () => {
  try {
    const res = await request.get('/api/user/info')
    if (res.data && res.data.id) {
      authStore.setUserInfo(res.data)
    }
  } catch (error) {
    console.warn('Failed to fetch user info', error)
  }
}

const fetchModels = async () => {
  try {
    const res = await request.get('/api/models/list')
    if (res.data && res.data.success) {
      availableModels.value = res.data.data.filter((m: any) => m.model_type === 4 || m.model_type === '4') // 只筛选图片生成模型
      if (availableModels.value.length > 0) {
        form.series_id = availableModels.value[0].series_id
        // set default values based on the first model
        if (currentModelRatios.value.length > 0) {
          form.aspect_ratio = currentModelRatios.value[0]
        }
        if (currentModelImageCounts.value.length > 0) {
          form.num_images = Number(currentModelImageCounts.value[0])
        }
      }
    }
  } catch (error) {
    // 忽略未授权等网络错误，仅打印警告
    console.warn('Failed to fetch models (this is normal if not logged in)', error)
  }
}

let contentMutationObserver: MutationObserver | null = null
let contentResizeObserver: ResizeObserver | null = null

const inputRef = ref<any | null>(null)

onMounted(() => {
  if (route.query.q) {
    if (route.path === '/explore') {
      router.push({ path: '/', query: { q: route.query.q } })
      return
    }
    form.prompt = route.query.q as string
    if (inputRef.value) {
      inputRef.value.focus()
    }
    router.replace({ path: route.path })
  }

  watch(() => route.query.q, (newQ) => {
    if (newQ) {
      if (route.path === '/explore') {
        router.push({ path: '/', query: { q: newQ } })
        return
      }
      form.prompt = newQ as string
      if (inputRef.value) {
        inputRef.value.focus()
      }
      router.replace({ path: route.path })
    }
  })

  window.addEventListener('resize', handleResize)
  document.addEventListener('click', handleClickOutside)
  
  if (contentAreaRef.value) {
    // 监听高度变化
    contentResizeObserver = new ResizeObserver(() => {
      if (isAutoScrolling.value && contentAreaRef.value) {
        contentAreaRef.value.scrollTop = contentAreaRef.value.scrollHeight
      }
    })
    
    contentResizeObserver.observe(contentAreaRef.value)
    Array.from(contentAreaRef.value.children).forEach(child => {
      contentResizeObserver!.observe(child)
    })

    // 监听 DOM 增加或改变
    contentMutationObserver = new MutationObserver((mutations) => {
      mutations.forEach(mutation => {
        mutation.addedNodes.forEach(node => {
          if (node.nodeType === 1) { // Element node
            contentResizeObserver?.observe(node as Element)
          }
        })
      })
      if (isAutoScrolling.value && contentAreaRef.value) {
        contentAreaRef.value.scrollTop = contentAreaRef.value.scrollHeight
      }
    })
    
    contentMutationObserver.observe(contentAreaRef.value, {
      childList: true,
      subtree: true,
      attributes: true,
      attributeFilter: ['src', 'style', 'class']
    })
  }

  timerInterval = window.setInterval(() => {
    currentTime.value = Date.now()
  }, 1000)
  fetchModels()
  fetchSettings()
  fetchPublicInspirations()
  const token = localStorage.getItem('token')
  if (token) {
  } else {
    // 首次访问或刷新时如果未登录，主动打开登录框
    authStore.openLogin()
    // 稍微延迟一下关闭骨架屏，有一个视觉上的“检查数据”的过渡效果
    setTimeout(() => {
      isInitialLoading.value = false
    }, 600)
  }
  const userInfoData = localStorage.getItem('userInfo')
  if (userInfoData) {
    try {
      const user = JSON.parse(userInfoData)
      authStore.setUserInfo(user)
    } catch (e) {}
  }
  if (isLoggedIn.value) {
    fetchUserInfo() // 立即获取最新用户信息（积分等）
    
    // 优先恢复保存的会话，以便立即开始加载任务，不等待 fetchSessions 完成
    const savedSessionId = sessionStorage.getItem('activeSessionId')
    if (savedSessionId) {
      isNewTaskMode.value = false
      activeSessionId.value = savedSessionId
      taskList.value = []
      scrollToBottom(true)
      fetchTasksBySession(savedSessionId)
      
      // 单独去获取这个会话的详细信息，如果它在一个被折叠的项目组或者第二页，我们需要展开/获取它
      request.get(`/api/user/sessions/detail?session_id=${savedSessionId}`).then(res => {
        if (res.data && res.data.code === 200 && res.data.data) {
          const session = res.data.data
          // 如果该会话属于某个项目组，并且项目组处于折叠状态，强制展开它
          if (session.project_id && session.project_id !== 0) {
            if (collapsedGroups.value.has(session.project_id)) {
              collapsedGroups.value.delete(session.project_id)
            }
            
            // 检查项目组是否存在于左侧列表中
            const existingGroup = projectGroups.value.find(g => g.id === session.project_id)
            if (!existingGroup) {
              // 项目组不在当前列表中（可能在很靠后的页数），单独拉取该项目组
              request.get(`/api/user/projects/detail?id=${session.project_id}`).then(groupRes => {
                if (groupRes.data && groupRes.data.code === 200 && groupRes.data.data) {
                  // 再次检查防止并发导致重复添加
                  if (!projectGroups.value.some(g => g.id === session.project_id)) {
                    projectGroups.value.push(groupRes.data.data)
                  }
                  // 顺便把这个项目组里的所有会话拉取一下，保证列表丰满
                  fetchSessionsByGroup(session.project_id)
                }
              })
            }

            // 我们还需要确保该会话在本地列表中，如果没有，把它推入对应的组中
            // 这样即使用户没有点击“加载更多项目组”，该会话也能显示在侧边栏
            const existingSession = sessionList.value.find(s => s.session_id === session.session_id)
            if (!existingSession) {
               sessionList.value.push(session)
            }
          } else {
             // 如果在最近任务中
             const existingSession = sessionList.value.find(s => s.session_id === session.session_id)
             if (!existingSession) {
               sessionList.value.push(session)
             }
          }
        }
      })
    }

    Promise.all([
      fetchProjectGroups().then(() => {
        // Fetch sessions for all groups and recent tasks
        const fetchPromises = [fetchSessionsByGroup(0)]
        projectGroups.value.forEach(g => {
          fetchPromises.push(fetchSessionsByGroup(g.id))
        })
        return Promise.all(fetchPromises)
      })
    ]).then(() => {
      if (activeSessionId.value) {
        // 页面刷新后，如果恢复了之前的会话，并且该会话在项目组里，确保它展开
        const session = sessionList.value.find(s => s.session_id === activeSessionId.value)
        if (session && session.project_id !== undefined && session.project_id !== 0) {
          if (collapsedGroups.value.has(session.project_id)) {
            collapsedGroups.value.delete(session.project_id)
          }
        }
      }
    })
    
    // 定时轮询获取最新状态（处理刷新后还有任务正在生成的情况）
    // 为了支持多窗口同步，移除 hasPending 条件限制，保持每5秒静默刷新
    pollInterval = window.setInterval(() => {
      if (isMutatingSession.value) return // 阻止在拖拽更新等过程中被定时刷新覆盖
      
      // 注意：静默刷新拉取全量（当前已加载页数 * 页面大小）数据，以便完全同步其他窗口的排序和加载更多操作
      fetchProjectGroups(false, true).then(() => {
        fetchSessionsByGroup(0, false, true)
        projectGroups.value.forEach(g => {
          fetchSessionsByGroup(g.id, false, true)
        })
      })
      fetchUserInfo()
      
      // 只要弹窗是开着的，就实时刷新积分明细记录
      if (showPointsRecordDialog.value) {
        fetchPointsRecords()
      }
      
      if (activeSessionId.value) {
        fetchTasksBySession(activeSessionId.value, false, true)
        
        // 同步当前激活会话的详细信息，防止在其他窗口被移动后本窗口不知情，导致强行保留在原错误分组中
        request.get(`/api/user/sessions/detail?session_id=${activeSessionId.value}`).then(res => {
          if (res.data && res.data.code === 200 && res.data.data) {
            const serverSession = res.data.data
            const existing = sessionList.value.find(s => s.session_id === activeSessionId.value)
            if (existing && (existing.project_id || 0) !== (serverSession.project_id || 0)) {
              existing.project_id = serverSession.project_id
              sessionList.value = [...sessionList.value]
              
              if (serverSession.project_id !== 0) {
                const groupExists = projectGroups.value.find(g => g.id === serverSession.project_id)
                if (!groupExists) {
                  request.get(`/api/user/projects/detail?id=${serverSession.project_id}`).then(groupRes => {
                    if (groupRes.data && groupRes.data.code === 200 && groupRes.data.data) {
                      if (!projectGroups.value.some(g => g.id === serverSession.project_id)) {
                        projectGroups.value.push(groupRes.data.data)
                      }
                      fetchSessionsByGroup(serverSession.project_id, false, true)
                    }
                  }).catch(() => {})
                }
              }
            }
          }
        }).catch(() => {})
      }
      
      // 对于后台正在生成，但不在当前聊天窗口的会话，我们也需要拉取它的任务以更新 tempTasksMap
      const backgroundSessions = Object.keys(tempTasksMap.value).filter(sid => 
        sid !== activeSessionId.value && tempTasksMap.value[sid].some((t: any) => t.status === 0)
      )
      backgroundSessions.forEach(sid => {
        fetchTasksBySession(sid, false, true)
      })
    }, 5000)
  }
})

onUnmounted(() => {
  if (contentMutationObserver) {
    contentMutationObserver.disconnect()
    contentMutationObserver = null
  }
  if (contentResizeObserver) {
    contentResizeObserver.disconnect()
    contentResizeObserver = null
  }
  window.removeEventListener('resize', handleResize)
  document.removeEventListener('click', handleClickOutside)
  if (pollInterval) {
    clearInterval(pollInterval)
  }
  if (timerInterval) {
    clearInterval(timerInterval)
  }
})

const groupPagination = ref<Record<number, { page: number, hasMore: boolean, loadingMore: boolean }>>({})
const initGroupPagination = (groupId: number) => {
  if (!groupPagination.value[groupId]) {
    groupPagination.value[groupId] = { page: 1, hasMore: false, loadingMore: false }
  }
}

const fetchSessionsByGroup = async (groupId: number, loadMore = false, silent = false) => {
  initGroupPagination(groupId)
  const pagination = groupPagination.value[groupId]

  if (loadMore && !pagination.hasMore) return
  if (loadMore) pagination.loadingMore = true

  try {
    const currentPage = pagination.page
    const pageToFetch = loadMore ? currentPage + 1 : 1
    
    let pageSize = 13
    if (!loadMore && currentPage > 1) {
      pageSize = currentPage * 13
    }
    
    const res = await request.get(`/api/user/sessions/list?project_id=${groupId}&page=${pageToFetch}&page_size=${pageSize}&_t=${Date.now()}`)
    
    if (res.data && res.data.list) {
      const list = res.data.list
      const total = res.data.total
      
      if (loadMore) {
        const newItems = list.filter((s: any) => !sessionList.value.some(exist => exist.session_id === s.session_id))
        sessionList.value.push(...newItems)
        pagination.page = pageToFetch
        
        if (list.length < 13 || total <= pageToFetch * 13) {
          pagination.hasMore = false
        }
      } else {
        // 不论是否静默刷新，只要是当前选中的会话，我们就得保住它！
        const currentActiveSession = sessionList.value.find(s => s.session_id === activeSessionId.value)
        
        if (!silent) {
          // 非静默刷新时，直接替换该组的列表
          sessionList.value = sessionList.value.filter(s => (s.project_id || 0) !== groupId)
          
          // 剔除即将在新列表中出现的会话，防止因为跨组移动导致的本地重复数据
          const newSessionIds = new Set(list.map((s: any) => s.session_id))
          sessionList.value = sessionList.value.filter(s => !newSessionIds.has(s.session_id))
          
          sessionList.value.push(...list)
          
          if (currentActiveSession && (currentActiveSession.project_id || 0) === groupId && !list.some((s: any) => s.session_id === currentActiveSession.session_id)) {
            sessionList.value.push(currentActiveSession)
          }
        } else {
          // 静默刷新时，直接使用服务器返回的列表和顺序，并保留不在当前组的和其他特殊会话
          let existingOtherSessions = sessionList.value.filter(s => (s.project_id || 0) !== groupId)
          
          // 剔除即将在当前组更新的会话，防止因为跨组移动导致的本地重复数据
          const newSessionIds = new Set(list.map((s: any) => s.session_id))
          existingOtherSessions = existingOtherSessions.filter(s => !newSessionIds.has(s.session_id))
          
          // 保留在 list 中的当前组 session，但是按照服务器的 list 顺序
          const updatedCurrentGroupSessions = list.map((serverSession: any) => {
            const existing = sessionList.value.find(s => s.session_id === serverSession.session_id)
            return existing ? { ...existing, ...serverSession } : serverSession
          })
          
          // 确保 activeSessionId 被保留，即使它不在当前拉取的列表中（可能在更深的页）
          if (currentActiveSession && (currentActiveSession.project_id || 0) === groupId && !list.some((s: any) => s.session_id === currentActiveSession.session_id)) {
            updatedCurrentGroupSessions.push(currentActiveSession)
          }
          
          sessionList.value = [...existingOtherSessions, ...updatedCurrentGroupSessions]
        }
        
        // 全局去重兜底，以防万一出现重复 key 导致渲染异常
        const uniqueSessions = []
        const seen = new Set()
        for (let i = sessionList.value.length - 1; i >= 0; i--) {
          const s = sessionList.value[i]
          if (!seen.has(s.session_id)) {
            seen.add(s.session_id)
            uniqueSessions.unshift(s)
          }
        }
        sessionList.value = uniqueSessions
        
        if (list.length < pageSize || total <= pageSize) {
          pagination.hasMore = false
        } else {
          pagination.hasMore = true
        }
      }
    }
  } catch (error: any) {
    console.error('Failed to fetch sessions for group', groupId, error)
  } finally {
    if (!silent && !loadMore && groupId === 0) isInitialLoading.value = false
    if (loadMore) pagination.loadingMore = false
  }
}

const hasMoreTasks = ref(false)
const currentTaskPage = ref(1)

const isTasksLoading = ref(false)
const isLoadingMoreTasks = ref(false)

const fetchTasksBySession = async (sessionId: string, loadMore = false, silent = false) => {
  if (loadMore && !hasMoreTasks.value) return

  if (!silent) {
    if (loadMore) {
      isLoadingMoreTasks.value = true
    } else {
      isTasksLoading.value = true
    }
  }

  try {
    const currentPage = currentTaskPage.value
    const pageToFetch = loadMore ? currentPage + 1 : 1
    
    let pageSize = 10
    if (!loadMore && currentPage > 1) {
      pageSize = currentPage * 10
    }
    
    const res = await request.get(`/api/user/tasks/list?session_id=${sessionId}&page=${pageToFetch}&page_size=${pageSize}`)
    
    if (res.data && res.data.list) {
      const dataList = res.data.list
      const total = res.data.total
      
      const newTasks = dataList.map((task: any) => ({
        ...task,
        num_images: task.num_images || (() => {
          if (task.image_url && task.image_url !== '[]') return parseImageUrls(task.image_url).length
          if (tempTasksMap.value[sessionId]) {
            const temp = tempTasksMap.value[sessionId].find(t => t.id === task.id)
            if (temp && temp.num_images) return temp.num_images
          }
          return 1
        })()
      }))
      
      // 无论后端的默认顺序是什么，我们在前端统一根据 created_at 进行严格的升序排序（旧的在上，新的在下）
      // 为了防止同一时间发送的多条消息乱序，如果时间相同则按 id 升序
      const sortedNewTasks = [...newTasks].sort((a: any, b: any) => {
        const timeA = new Date(a.created_at).getTime()
        const timeB = new Date(b.created_at).getTime()
        if (timeA === timeB) {
          if (typeof a.id === 'number' && typeof b.id === 'number') {
            return a.id - b.id
          }
          return String(a.id).localeCompare(String(b.id))
        }
        return timeA - timeB
      })

      // 如果当前拉取的不是正在查看的会话（后台会话），我们只更新它的临时状态，不影响 taskList
      if (sessionId !== activeSessionId.value) {
        if (tempTasksMap.value[sessionId]) {
          sortedNewTasks.forEach((newTask: any) => {
            const tempIndex = tempTasksMap.value[sessionId].findIndex(t => t.id === newTask.id)
            if (tempIndex > -1) {
              if (newTask.status !== 0) {
                // 后台任务已完成，从临时队列中移除
                tempTasksMap.value[sessionId].splice(tempIndex, 1)
              } else {
                tempTasksMap.value[sessionId][tempIndex].num_images = newTask.num_images
              }
            }
          })
          // 如果临时队列空了，直接删除该键
          if (tempTasksMap.value[sessionId].length === 0) {
            delete tempTasksMap.value[sessionId]
          }
        }
        return
      }
      
      if (loadMore) {
        // 保存当前的滚动高度和位置
        let oldScrollHeight = 0
        let oldScrollTop = 0
        let anchorId: string | null = null
        let anchorOffsetTop = 0
        
        if (contentAreaRef.value) {
          oldScrollHeight = contentAreaRef.value.scrollHeight
          oldScrollTop = contentAreaRef.value.scrollTop
          
          // 尝试找到当前的第一个任务元素作为锚点，这样恢复滚动最精确
          const firstOldTask = taskList.value.length > 0 ? taskList.value[0] : null
          if (firstOldTask) {
            anchorId = 'task-' + firstOldTask.id
            const anchorEl = document.getElementById(anchorId)
            if (anchorEl) {
              anchorOffsetTop = anchorEl.offsetTop
            }
          }
        }

        // 当向上滚动加载更旧的消息时，这些更旧的消息（已经排好序）应该拼接在当前列表的最前面
        const newItems = sortedNewTasks.filter((t: any) => !taskList.value.some(exist => exist.id === t.id))
        taskList.value = [...newItems, ...taskList.value]
        currentTaskPage.value = pageToFetch
        
        if (dataList.length < pageSize || total <= pageToFetch * pageSize) {
          hasMoreTasks.value = false
        }

        // 等待 Vue 更新 DOM 后，恢复滚动位置
        nextTick(() => {
          if (contentAreaRef.value) {
            const el = contentAreaRef.value
            
            if (anchorId) {
              const anchorEl = document.getElementById(anchorId)
              if (anchorEl) {
                // 恢复锚点元素相对于视口的相对位置
                const relativeTop = anchorOffsetTop - oldScrollTop
                el.scrollTop = anchorEl.offsetTop - relativeTop
              } else {
                // 如果没找到锚点，降级使用高度差计算
                const newScrollHeight = el.scrollHeight
                el.scrollTop = newScrollHeight - oldScrollHeight + oldScrollTop
              }
            } else {
              const newScrollHeight = el.scrollHeight
              el.scrollTop = newScrollHeight - oldScrollHeight + oldScrollTop
            }
          }
        })
      } else {
        if (!silent) {
          // 首次加载或刷新整个页面，直接赋值（注意此时已经按时间升序排好）
          taskList.value = sortedNewTasks
          
          if (tempTasksMap.value[sessionId]) {
             tempTasksMap.value[sessionId] = tempTasksMap.value[sessionId].filter(t => 
               t.status === 0 && !newTasks.some((nt: any) => nt.id === t.id)
             )
          }
          
          // 如果是首次加载页面，立刻跳转到底部（不带平滑动画），避免因为 DOM 渲染导致的高度变化引发闪烁跳动
          scrollToBottom(true)
          
          if (dataList.length < pageSize || total <= pageSize) {
            hasMoreTasks.value = false
          } else {
            hasMoreTasks.value = true
          }
        } else {
          // silent fetch logic
          let hasNew = false
          // 对于静默刷新，我们要按时间顺序（从早到晚）遍历
          sortedNewTasks.forEach((newTask: any) => {
            const existTask = taskList.value.find(t => t.id === newTask.id)
            if (existTask) {
              if (existTask.status !== newTask.status || existTask.image_url !== newTask.image_url || existTask.log_content !== newTask.log_content) {
                existTask.status = newTask.status
                existTask.image_url = newTask.image_url
                existTask.log_content = newTask.log_content
                existTask.num_images = newTask.num_images
                // 给渲染一点时间再滚动
                setTimeout(() => {
                  scrollToBottom()
                }, 100)
              }
            } else {
              let inTemp = false
              if (tempTasksMap.value[sessionId]) {
                const tempIndex = tempTasksMap.value[sessionId].findIndex(t => t.id === newTask.id)
                if (tempIndex > -1) {
                  inTemp = true
                  if (newTask.status !== 0 || newTask.image_url) {
                    tempTasksMap.value[sessionId].splice(tempIndex, 1)
                    // 如果原先在临时队列里，说明是最新生成的任务，我们 push 到列表末尾
                    taskList.value.push(newTask)
                    hasNew = true
                  } else {
                    tempTasksMap.value[sessionId][tempIndex].num_images = newTask.num_images
                  }
                }
              }
              if (!inTemp) {
                // 如果不仅不在列表里，也不在临时队列里，说明也是新任务，同样 push 到末尾
                taskList.value.push(newTask)
                hasNew = true
              }
            }
          })
          if (hasNew) {
            // 给渲染一点时间再滚动
            setTimeout(() => {
              scrollToBottom()
            }, 100)
          }
          
          if (dataList.length < pageSize || total <= pageSize) {
            hasMoreTasks.value = false
          } else {
            hasMoreTasks.value = true
          }
        }
      }
    } else {
      if (!silent) {
        taskList.value = []
        hasMoreTasks.value = false
      }
    }
  } catch (error: any) {
    console.error('Failed to fetch tasks', error)
  } finally {
    if (!silent) {
      if (loadMore) {
        isLoadingMoreTasks.value = false
      } else {
        isTasksLoading.value = false
      }
    }
  }
}

const selectSession = (sessionId: string) => {
  if (activeSessionId.value === sessionId) return
  
  deletedTaskIds.value.clear()
  isNewTaskMode.value = false
  activeSessionId.value = sessionId
  sessionStorage.setItem('activeSessionId', sessionId)
  
  // 在切换会话时，先清空当前的任务列表，以便触发骨架屏
  taskList.value = []
  currentTaskPage.value = 1
  // 去掉这里切换会话时的强制滚动到底部，避免跳动
  // scrollToBottom(true)
  
  // 如果选中的会话在某个项目组中，确保该项目组是展开状态
  const session = sessionList.value.find(s => s.session_id === sessionId)
  if (session && session.project_id !== undefined && session.project_id !== 0) {
    if (collapsedGroups.value.has(session.project_id)) {
      collapsedGroups.value.delete(session.project_id)
    }
    
    // 如果这个项目组目前不在可见的项目组列表中，我们需要单独拉取它并插入列表，
    // 以确保用户能看到它，同时不影响原本的分页逻辑。
    const groupExists = projectGroups.value.find(g => g.id === session.project_id)
    if (!groupExists) {
      request.get(`/api/user/projects/detail?id=${session.project_id}`).then(groupRes => {
        if (groupRes.data && groupRes.data.code === 200 && groupRes.data.data) {
          if (!projectGroups.value.some(g => g.id === session.project_id)) {
            // 将其放在列表顶部或合适位置
            projectGroups.value.unshift(groupRes.data.data)
          }
          fetchSessionsByGroup(session.project_id, false, true)
        }
      }).catch(() => {})
    }
  }
  
  // 不再清空 tempTasksMap，这样切换回来时还能看到临时任务
  fetchTasksBySession(sessionId)
  // 如果在移动端，选择会话后自动关闭侧边栏
  if (isMobileMenuOpen.value) {
    isMobileMenuOpen.value = false
  }
}

const resetTask = () => {
  if (!isLoggedIn.value) {
    authStore.openLogin()
    return
  }
  activeSessionId.value = null
  sessionStorage.removeItem('activeSessionId')
  taskList.value = []
  currentTaskPage.value = 1
  form.prompt = ''
  clearReferenceImage()
  isNewTaskMode.value = true
  hasMoreTasks.value = false
  isMobileMenuOpen.value = false
}

const handleNewChatFromExplore = () => {
  resetTask()
  router.push('/')
}

const scrollToBottom = (force = false) => {
  if (force) {
    isAutoScrolling.value = true
  }
  nextTick(() => {
    if (contentAreaRef.value) {
      if (force || isAutoScrolling.value) {
        contentAreaRef.value.scrollTop = contentAreaRef.value.scrollHeight
      }
    }
  })
}

const handleEnter = (e: KeyboardEvent) => {
  if (e.isComposing || e.shiftKey) return
  handleGenerate()
}

const regenerateTask = (task: any) => {
  form.prompt = task.prompt
  form.series_id = task.series_id || 'default'
  form.resolution = task.resolution || '1K'
  form.aspect_ratio = task.size || 'auto'
  form.num_images = task.num_images || 1
  try {
    form.reference_image = task.reference_image ? JSON.parse(task.reference_image) : []
  } catch (e) {
    form.reference_image = task.reference_image ? [task.reference_image] : []
  }
  pendingReferenceFiles.value = []
  handleGenerate()
}

const reEditTask = (task: any) => {
  form.prompt = task.prompt
  form.series_id = task.series_id || 'default'
  form.resolution = task.resolution || '1K'
  form.aspect_ratio = task.size || 'auto'
  form.num_images = task.num_images || 1
  try {
    form.reference_image = task.reference_image ? JSON.parse(task.reference_image) : []
  } catch (e) {
    form.reference_image = task.reference_image ? [task.reference_image] : []
  }
  pendingReferenceFiles.value = []
  isInputFocused.value = true
  nextTick(() => {
    const inputEl = document.querySelector('.main-input textarea') as HTMLTextAreaElement
    if (inputEl) {
      inputEl.focus()
    }
  })
}

const handleTaskCommand = async (command: string, task: any) => {
  if (command === 'delete') {
    try {
      await ElMessageBox.confirm(
        '删除的历史记录将无法找回，确认删除该批次结果吗？',
        '删除确认',
        {
          confirmButtonText: '确认删除',
          cancelButtonText: '取消',
          type: 'warning',
          confirmButtonClass: 'el-button--danger'
        }
      )
      
      const res = await request.post('/api/user/tasks/delete', {
        task_id: task.id
      })
      
      if (res.data.code === 200) {
        ElMessage.success('删除成功')
        deletedTaskIds.value.add(task.id)
        // Remove from taskList using re-assignment for sure reactivity, and loose equality for ID
        taskList.value = taskList.value.filter(t => t.id != task.id)
        
        // Also remove from tempTasksMap if exists
        if (activeSessionId.value && tempTasksMap.value[activeSessionId.value]) {
          tempTasksMap.value[activeSessionId.value] = tempTasksMap.value[activeSessionId.value].filter(t => t.id != task.id)
        }
      } else {
        ElMessage.error(res.data.message || '删除失败')
      }
    } catch (e) {
      if (e !== 'cancel') {
        ElMessage.error('删除失败')
      }
    }
  }
}

const handleGenerate = async () => {
  if (!isLoggedIn.value) {
    authStore.openLogin()
    return
  }
  
  if (!form.prompt.trim()) {
    return
  }

  isNewTaskMode.value = false
  
  let finalReferenceImages = [...form.reference_image]
  if (pendingReferenceFiles.value.length > 0) {
    for (let i = 0; i < pendingReferenceFiles.value.length; i++) {
      const file = pendingReferenceFiles.value[i]
      const formData = new FormData()
      formData.append('file', file)
      try {
        const res = await request.post('/api/upload', formData, {
          headers: { 'Content-Type': 'multipart/form-data' }
        })
        if (res.data && res.data.code === 200) {
          finalReferenceImages[i] = res.data.data.url
        } else {
          ElMessage.error(res.data.message || '上传参考图失败')
          return
        }
      } catch (e) {
        ElMessage.error('上传参考图失败')
        return
      }
    }
  }

  const currentPrompt = form.prompt
  const currentResolution = form.resolution
  const currentAspectRatio = form.aspect_ratio
  const currentNumImages = form.num_images
  const currentReferenceImage = finalReferenceImages.length > 0 ? JSON.stringify(finalReferenceImages) : ''
  const currentSize = currentAspectRatio || "1:1"
  const currentSeriesId = form.series_id
  
  // 提前记录是否从探索页面过来
  const isFromExplore = route.path === '/explore'

  // 清空输入和参考图，跳转路由（确保跳转不阻断后续流程）
  form.prompt = ''
  clearReferenceImage() // clear after send
  if (isFromExplore) {
    router.push('/')
  }
  
  if (!activeSessionId.value || isFromExplore) {
    activeSessionId.value = `session_${Date.now()}`
    sessionStorage.setItem('activeSessionId', activeSessionId.value)
    // 立即在侧边栏显示新建的任务会话
    sessionList.value.unshift({
      session_id: activeSessionId.value,
      title: currentPrompt,
      image_url: '',
      status: 0, // 初始状态标记为生成中
      project_id: 0
    })
    hasMoreTasks.value = false
    // 如果是从探索页过来的新对话，清空当前页面的旧任务列表
    if (isFromExplore) {
      taskList.value = []
    }
  } else {
    // 如果已有当前会话正在生成新内容，确保该会话在侧边栏最上方
    const index = sessionList.value.findIndex(s => s.session_id === activeSessionId.value)
    if (index > 0) {
      const session = sessionList.value[index]
      session.status = 0 // 将该会话状态标记为生成中，以便显示 loading
      session.title = currentPrompt // 更新为最新的 prompt
      sessionList.value.splice(index, 1)
      sessionList.value.unshift(session)
    } else if (index === 0) {
      sessionList.value[0].status = 0
      sessionList.value[0].title = currentPrompt
    }
  }

  const currentSessionId = activeSessionId.value
  
  // Find current session to get its project_id
  let currentProjectId = 0
  if (currentSessionId) {
    const session = sessionList.value.find(s => s.session_id === currentSessionId)
    if (session && session.project_id !== undefined) {
      currentProjectId = session.project_id
      // 如果当前会话在某个项目组中，确保该项目组是展开状态
      if (currentProjectId !== 0 && collapsedGroups.value.has(currentProjectId)) {
        collapsedGroups.value.delete(currentProjectId)
      }
    }
  }
  
  if (currentProjectId === 0) {
    isRecentTasksExpanded.value = true
  } else {
    isProjectGroupsExpanded.value = true
  }

  const tempId = Date.now()
  const tempTask = {
    id: tempId,
    prompt: currentPrompt,
    size: currentSize,
    resolution: currentResolution,
    series_id: currentSeriesId,
    reference_image: currentReferenceImage,
    status: 0,
    image_url: '',
    log_content: '',
    num_images: currentNumImages || 1,
    created_at: new Date().toISOString()
  }
  
  if (!tempTasksMap.value[currentSessionId]) {
    tempTasksMap.value[currentSessionId] = []
  }
  tempTasksMap.value[currentSessionId].push(tempTask)
  // 给渲染一点时间再滚动
  setTimeout(() => {
    scrollToBottom(true)
  }, 100)

  try {
    // 现在的后端接口 /api/generate 已经改成了异步，会立刻返回 200
    const response = await request.post(`/api/generate`, {
      prompt: currentPrompt,
      resolution: currentResolution,
      aspect_ratio: currentAspectRatio,
      num_images: currentNumImages,
      reference_image: currentReferenceImage,
      series_id: currentSeriesId,
      session_id: currentSessionId,
      project_id: currentProjectId
    })
    
    if (response.data.code === 200) {
      // 成功提交给后端，此时让后台去慢慢生成。
      // 我们把前端临时生成的 id 替换成后端返回的真实 task_id，这样后续轮询能对应上
      if (response.data.data && response.data.data.task_id) {
        const target = tempTasksMap.value[currentSessionId].find(t => t.id === tempId)
        if (target) {
          target.id = response.data.data.task_id
        }
      }
      // 不在这里弹"图片生成成功"，因为这只是任务创建成功，真正的图片还没回来
      // 也不需要去调用 fetchTasksBySession，交给那个每 5 秒的 pollInterval 去自动发现并刷新即可
      fetchUserInfo() // 立即更新用户积分
    } else {
      if (tempTasksMap.value[currentSessionId]) {
        const target = tempTasksMap.value[currentSessionId].find(t => t.id === tempId)
        if (target) {
          target.status = 2
          target.log_content = response.data.message || '生成失败'
        }
      }
    }
  } catch (error: any) {
    console.error(error)
    const errorMsg = error.response?.data?.message || error.response?.data?.error || '请求失败，请检查网络或后端服务'
    const returnedTaskId = error.response?.data?.data?.task_id
    if (tempTasksMap.value[currentSessionId]) {
      const target = tempTasksMap.value[currentSessionId].find(t => t.id === tempId)
      if (target) {
        if (returnedTaskId) {
          target.id = returnedTaskId
        }
        target.status = 2
        target.log_content = errorMsg
      }
    }
  } finally {
    if (activeSessionId.value === currentSessionId) {
      scrollToBottom(true)
    }
  }
}

const getAspectRatioStyle = (size: string) => {
  if (!size) return { aspectRatio: '1 / 1' }
  if (size.toLowerCase() === 'auto') return { aspectRatio: '1 / 1' }
  
  // 处理 1024x1024 这种格式
  if (size.includes('x') || size.includes('*')) {
    const parts = size.toLowerCase().split(/[x*]/)
    if (parts.length === 2 && !isNaN(Number(parts[0])) && !isNaN(Number(parts[1]))) {
      return { aspectRatio: `${parts[0]} / ${parts[1]}` }
    }
  }
  
  return { aspectRatio: size.replace(':', ' / ') }
}

const parseImageUrls = (imgData: string) => {
  if (!imgData) return []
  try {
    const arr = JSON.parse(imgData)
    if (Array.isArray(arr)) {
      // 对图片 URL 进行排序，确保多次加载或乱序返回时，同一批图片的展示顺序始终一致
      return arr.sort()
    }
  } catch (e) {
    // 兼容旧数据
  }
  return [imgData]
}

const getErrorMessage = (logContent?: string, status?: number) => {
  if (!logContent) return '未知错误'
  
  // 1. 尝试从日志中提取出上游 API 返回的错误 JSON 并获取其中的 message
  try {
    const errorBodyRegex = /Upstream API returned error status.*?\nBody:\s*({.*})/i;
    const match = logContent.match(errorBodyRegex);
    if (match && match[1]) {
      const errorJson = JSON.parse(match[1]);
      if (errorJson?.error?.message) {
        // 如果提取到了 message，根据错误内容做简单翻译（可选），或者直接抛出
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
    // JSON 解析失败则继续走下面的正则兜底逻辑
  }

  // 2. 如果没提取到标准的 JSON message，降级走之前的关键字匹配逻辑
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
  
  if (safeText.includes('timeout') || safeText.includes('network')) {
    return status === 3 ? '部分图片生成超时' : '生成超时，请检查网络或稍后重试'
  }
  
  // 兜底提示，避免暴露内部报错详情
  return status === 3 ? '部分图片生成失败' : '生成失败，请稍后重试'
}

const formatImageUrl = (img: string) => {
  if (!img) return ''

  let urlStr = img;

  // Extract array if it's a JSON string
  try {
    const arr = JSON.parse(img)
    if (Array.isArray(arr) && arr.length > 0) {
      urlStr = arr[0]
    }
  } catch (e) {
    // not json
  }

  if (urlStr.length > 1000 && !urlStr.startsWith('http') && !urlStr.startsWith('data:')) {
    return 'data:image/png;base64,' + urlStr
  }

  // Strip any existing backend origin to make it relative
  try {
    if (urlStr.startsWith('http')) {
      const urlObj = new URL(urlStr)
      if (urlObj.pathname.startsWith('/uploads/')) {
        urlStr = urlObj.pathname
      }
    }
  } catch(e) {}

  if (urlStr.startsWith('/uploads/')) {
    let defaultBase = ''
    if ((window as any).APP_CONFIG?.API_BASE_URL) {
      defaultBase = (window as any).APP_CONFIG.API_BASE_URL.replace(/\/$/, '')
    } else {
      // 动态获取当前访问的域名，并默认后端在 8088 端口
      // 如果你的前端和后端在线上是同域名、同端口（例如都通过 Nginx 代理到 80），可以直接用 window.location.origin
      if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
        defaultBase = window.location.protocol + "//" + window.location.hostname + ":8088"
      } else {
        defaultBase = window.location.origin // 线上环境默认使用当前访问的域名
      }
    }
    
    const base = imageDomain.value || defaultBase
    // Ensure no double slash
    return base.replace(/\/$/, '') + urlStr
  }

  return urlStr
}

const handleLogout = () => {
  authStore.logout()
  sessionStorage.removeItem('activeSessionId')
  sessionList.value = []
  taskList.value = []
  tempTasksMap.value = {}
  projectGroups.value = []
  activeSessionId.value = null
  isNewTaskMode.value = true
  hasMoreTasks.value = false
  clearReferenceImage()
  ElMessage.success('已退出登录')
  authStore.openLogin()
}

const showPointsRecordDialog = ref(false)
const isPointsVisible = ref(true)
const showRedeemDialog = ref(false)
const pointsRecordLoading = ref(false)
const cdkeyValue = ref('')
const redeemingCdkey = ref(false)
const pointsRecordList = ref<any[]>([])
const pointsRecordPage = ref(1)
const pointsRecordPageSize = ref(10)
const pointsRecordTotal = ref(0)

const fetchPointsRecords = async () => {
  pointsRecordLoading.value = true
  try {
    const res = await request.get(`/api/user/points/records?page=${pointsRecordPage.value}&page_size=${pointsRecordPageSize.value}`)
    if (res.data && res.data.list) {
      pointsRecordList.value = res.data.list
      pointsRecordTotal.value = res.data.total || 0
    }
  } catch (error) {
    ElMessage.error('获取积分记录失败')
  } finally {
    pointsRecordLoading.value = false
  }
}

const handleRedeemCdkey = async () => {
  if (!cdkeyValue.value.trim()) {
    ElMessage.warning('请输入兑换码')
    return
  }
  
  redeemingCdkey.value = true
  try {
    const res = await request.post('/api/user/cdkeys/use', { cdkey: cdkeyValue.value.trim() })
    ElMessage.success(res.data.message)
    cdkeyValue.value = ''
    
    // 更新用户信息和积分明细
    await fetchUserInfo()
    pointsRecordPage.value = 1
    fetchPointsRecords()
    
    // 关闭兑换弹窗
    showRedeemDialog.value = false
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '兑换失败')
  } finally {
    redeemingCdkey.value = false
  }
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

const handleCommand = (command: string) => {
  if (command === 'logout') {
    handleLogout()
  } else if (command === 'pointsRecord') {
    pointsRecordPage.value = 1
    showPointsRecordDialog.value = true
    fetchPointsRecords()
  } else if (command === 'recharge') {
    handleRecharge()
  }
}

// 只要积分明细弹窗是开着的，就会通过上面 2197 行的定时器实时刷新数据
watch(showPointsRecordDialog, (newVal) => {
  if (newVal) {
    fetchUserInfo()
    fetchPointsRecords()
  }
})

const cachedSettings = localStorage.getItem('site_settings')
const siteSettings = ref(cachedSettings ? JSON.parse(cachedSettings) : {})

// 监听设置变化
window.addEventListener('storage', (e) => {
  if (e.key === 'site_settings' && e.newValue) {
    siteSettings.value = JSON.parse(e.newValue)
  }
})

const handleRecharge = () => {
  const link = siteSettings.value?.recharge_link
  if (link) {
    window.open(link, '_blank')
  } else {
    ElMessage.info('充值链接未配置，请联系管理员')
  }
}

const useSuggestion = (text: string) => {
  form.prompt = text
}

const downloadImage = (url: string) => {
  if (!url) return
  const a = document.createElement('a')
  a.href = url
  a.download = `generated_image_${Date.now()}.png`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}
</script>

<style scoped>
* {
  -webkit-tap-highlight-color: transparent;
}

.empty-session-screen {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;
  padding-top: 100px;
}

.custom-empty-icon {
  position: relative;
  width: 100px;
  height: 100px;
  margin: 0 auto;
}

.icon-layer {
  position: absolute;
  inset: 0;
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.layer-1 {
  background: #f1f5f9;
  transform: rotate(-10deg) scale(0.9);
  animation: pulse-layer1 4s ease-in-out infinite alternate;
}

.layer-2 {
  background: #e2e8f0;
  transform: rotate(10deg) scale(0.95);
  opacity: 0.6;
  animation: pulse-layer2 4s ease-in-out infinite alternate-reverse;
}

.layer-3 {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 12px 24px rgba(15, 23, 42, 0.04);
  transform: rotate(0);
}

.layer-3 .el-icon {
  font-size: 42px;
  color: #334155;
  transition: transform 0.3s ease;
}

.custom-empty-icon:hover .layer-3 .el-icon {
  transform: scale(1.1) rotate(-5deg);
}

@keyframes pulse-layer1 {
  0% { transform: rotate(-10deg) scale(0.9); }
  100% { transform: rotate(-15deg) scale(0.95); }
}

@keyframes pulse-layer2 {
  0% { transform: rotate(10deg) scale(0.95); }
  100% { transform: rotate(15deg) scale(1); }
}

/* 全局布局重置与字体设置已在 App.vue，这里仅限定作用域 */
.layout-container {
    display: flex;
    height: 100dvh;
    width: 100%;
    position: relative;
    background-color: #ffffff;
    overflow: hidden;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  }

/* Mini Sidebar Styles */
.mini-sidebar {
    width: 64px;
    background-color: #ffffff;
    border-right: 1px solid #f3f4f6;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: flex-start;
    padding: 16px 0 72px 0;
    flex-shrink: 0;
    z-index: 30;
    box-sizing: border-box;
  }

.mini-sidebar-top {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
    flex: 1;
    overflow-y: auto;
    min-height: 0;
  }
  
  .mini-sidebar-top::-webkit-scrollbar {
    display: none;
  }
  
  .logo-mini {
    margin-bottom: 24px;
    cursor: pointer;
    flex-shrink: 0;
  }
  
  .logo-img {
    width: 32px;
    height: 32px;
    border-radius: 8px;
  }
  
  .nav-items {
    display: flex;
    flex-direction: column;
    gap: 12px;
    width: 100%;
    align-items: center;
    flex-shrink: 0;
  }

.nav-item {
  width: 48px;
  height: 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  cursor: pointer;
  color: #6b7280;
  transition: all 0.2s;
  gap: 4px;
}

.nav-item .el-icon {
  font-size: 20px;
}

.nav-item span {
  font-size: 10px;
}

.nav-item:hover {
  background-color: #f3f4f6;
  color: #111827;
}

.nav-item.active {
  background-color: #f3f4f6;
  color: #111827;
  font-weight: 600;
}

.mini-sidebar-bottom {
    display: flex;
    flex-direction: column;
    align-items: center;
    flex-shrink: 0;
    margin-top: 16px;
  }

.user-profile-mini {
  cursor: pointer;
  transition: opacity 0.2s;
}

.user-profile-mini:hover {
  opacity: 0.8;
}

/* History Sidebar Styles */
.history-sidebar {
    width: 260px;
    background-color: #f9fafb;
    border-right: 1px solid #f3f4f6;
    display: flex;
    flex-direction: column;
    flex-shrink: 0;
    z-index: 20;
    position: relative; /* 确保 z-index 正常工作 */
  }

  .global-account-container {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 324px; /* 64px (mini-sidebar) + 260px (history-sidebar) */
    height: 64px;
    z-index: 2000; /* 确保层级高于 ExploreContent 及其他滚动内容 */
    background: #f9fafb; /* 统一底色 */
    border-top: 1px solid #f3f4f6;
    border-right: 1px solid #f3f4f6;
    box-sizing: border-box; /* 防止边框撑大元素 */
    transition: width 0.3s ease;
  }

  .global-account-container.explore-mode {
    width: 64px;
    border-right: none;
    border-top: none;
    background: transparent;
  }
  
  .account-dropdown {
    width: 100%;
    height: 100%;
  }

  @media screen and (max-width: 768px) {
    .global-account-container {
      width: 100%;
      border-right: none;
      transform: translateX(-100%);
      transition: transform 0.3s ease;
      z-index: 2000;
    }
    
    .global-account-container.mobile-open {
      transform: translateX(0);
    }
  }

  .account-dropdown {
    width: 100%;
    height: 100%;
  }

  .accountTrigger {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0;
    transition: background-color 0.2s ease;
  }

  .global-account-container.explore-mode .accountTrigger {
    justify-content: center;
  }

  .accountTrigger:hover {
    background-color: #f3f4f6;
  }

  .accountTriggerAvatar {
    width: 64px;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    background: transparent; /* 移除原来的白色背景，使其统一 */
  }

  .accountTriggerInfo {
    flex: 1;
    height: 100%;
    display: flex;
    align-items: center;
    padding-right: 16px;
    overflow: hidden;
    background: transparent; /* 移除原来的底色，使其统一 */
  }

  .accountTriggerNameWrap {
    display: flex;
    flex-direction: column;
    flex: 1;
    overflow: hidden;
    justify-content: center;
    align-items: flex-start;
    padding-right: 8px;
  }

  .accountTriggerName {
    font-size: 14px;
    font-weight: 500;
    color: #111827;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: left;
    width: 100%;
    line-height: 1.2;
  }

  .accountTriggerEmail {
    font-size: 12px;
    color: #6b7280;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    text-align: left;
    width: 100%;
    margin-top: 2px;
    line-height: 1.2;
  }

  .accountTriggerMembership {
    flex-shrink: 0;
    display: flex;
    align-items: center;
  }

  .accountHostTag {
  display: inline-flex;
  align-items: center;
  font-size: 12px;
  padding: 2px 6px;
  background-color: #f3f4f6;
  color: #6b7280;
  border-radius: 4px;
  border: 1px solid #e5e7eb;
}

  .login-tag {
    background-color: #eff6ff;
    color: #3b82f6;
    border-color: #bfdbfe;
  }

  .history-header {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
    flex-shrink: 0;
  }

.history-header-title {
  font-weight: 600;
  font-size: 16px;
  color: #111827;
}

.new-chat-icon {
  font-size: 40px !important;
  color: #111827;
  cursor: pointer;
  padding: 12px;
  border-radius: 8px;
  transition: background-color 0.2s ease;
}

.new-chat-icon:hover {
  background-color: #e5e7eb;
}

.sidebar-content {
  flex: 1;
  padding: 8px;
  padding-bottom: 72px; /* 留出底部用户信息空间 */
  overflow-y: auto;
  position: relative;
  z-index: 10;
}

/* 隐藏滚动条 */
.sidebar-content::-webkit-scrollbar {
  width: 4px;
}
.sidebar-content::-webkit-scrollbar-thumb {
  background: transparent;
  border-radius: 4px;
}
.sidebar:hover .sidebar-content::-webkit-scrollbar-thumb {
  background: #d1d5db;
}

.history-list {
  margin-top: 8px;
}

.history-title {
  font-size: 12px;
  color: #888888;
  margin-bottom: 8px;
  padding-left: 12px;
  font-weight: 600;
}

.project-groups {
  margin-bottom: 16px;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 12px;
  margin-bottom: 8px;
  -webkit-tap-highlight-color: transparent;
}

.group-title {
  font-size: 12px;
  color: #888888;
  font-weight: 600;
}

.add-group-icon {
  color: #6b7280;
  cursor: pointer;
  border-radius: 8px;
  display: flex !important;
  align-items: center;
  justify-content: center;
}

.add-group-icon svg {
  width: 100% !important;
  height: 100% !important;
}

.add-group-icon:hover {
  background-color: #e5e7eb;
  color: #111827;
}

.project-group-container {
  margin-bottom: 8px;
}

.group-session-list {
  padding-left: 24px; /* 缩进以体现层级关系 */
}

.project-group-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 8px;
  cursor: pointer;
  color: #374151;
  font-size: 14px;
  transition: all 0.2s ease;
  margin-bottom: 2px;
  -webkit-tap-highlight-color: transparent;
}

.project-group-item:hover {
  background-color: transparent;
}

.group-folder-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 6px;
  background-color: #f3f4f6;
  border: 1px solid #e5e7eb;
  transition: all 0.2s ease;
  margin-right: 8px;
}

.project-group-item .el-icon {
  font-size: 16px;
  color: #6b7280;
}

.project-group-item:hover .group-folder-icon {
  background-color: #e5e7eb;
}

.project-group-item.active {
  background-color: #f3f4f6;
  color: #111827;
  font-weight: 500;
}

.project-group-item.active .group-folder-icon {
  background-color: #eff6ff;
  border-color: #bfdbfe;
}

.project-group-item.active .group-folder-icon .el-icon {
  color: #3b82f6;
}

.group-name {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}

.group-content {
  flex: 1;
  display: flex;
  align-items: center;
  overflow: hidden;
  min-width: 0;
}

.group-actions {
  opacity: 0;
  pointer-events: none;
  flex-shrink: 0;
}

.project-group-item:hover .group-actions,
.project-group-item.dropdown-open .group-actions {
  opacity: 1;
  pointer-events: auto;
}

.history-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  color: #0d0d0d;
  transition: background-color 0.2s ease;
  font-size: 14px;
  margin-bottom: 2px;
  position: relative;
  overflow: hidden; /* 防止内部内容超出 */
  -webkit-tap-highlight-color: transparent;
}

.history-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0; /* 允许内容收缩 */
  overflow: hidden;
}

.history-text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  min-width: 0;
}

.history-item .el-icon {
  font-size: 16px;
  color: #0d0d0d;
}

.history-thumbnail {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.thumbnail-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 4px;
}

.history-item:hover {
  background-color: #e5e7eb;
}

.history-item.active {
  background-color: #e5e7eb;
  font-weight: 500;
}

.history-content {
  flex: 1;
  display: flex;
  align-items: center;
  overflow: hidden;
}

.history-text {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  flex: 1;
  min-width: 0;
}

.history-actions {
  opacity: 0;
  pointer-events: none;
  flex-shrink: 0;
}

/* Sidebar Skeleton Styles */
.sidebar-skeleton {
  padding: 0 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.skeleton-history-item {
  display: flex;
  align-items: center;
  padding: 8px;
  border-radius: 8px;
  background-color: transparent;
}

.skeleton-history-thumbnail {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  background-color: #f3f4f6;
  margin-right: 12px;
  flex-shrink: 0;
  animation: pulse 1.5s infinite;
}

.skeleton-history-text {
  height: 16px;
  background-color: #f3f4f6;
  border-radius: 4px;
  flex: 1;
  max-width: 80%;
  animation: pulse 1.5s infinite;
}

@keyframes pulse {
  0% { opacity: 0.6; }
  50% { opacity: 1; }
  100% { opacity: 0.6; }
}

.history-item:hover .history-actions,
.history-item.dropdown-open .history-actions {
  opacity: 1;
  pointer-events: auto;
}

.action-icon {
  font-size: 32px !important;
  padding: 10px;
  border-radius: 8px;
  color: #6b7280;
}

.action-icon:hover {
  background-color: #d1d5db;
  color: #111827;
}

.text-danger {
  color: #f56c6c;
}

.no-history {
  padding: 16px 12px;
  text-align: center;
  color: #888888;
  font-size: 13px;
}

.load-more {
  padding: 8px 12px;
  margin: 4px 12px 12px 12px;
  text-align: center;
  color: #9ca3af;
  font-size: 12px;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  background-color: transparent;
}

.load-more:hover {
  background-color: #f3f4f6;
  color: #4b5563;
}

/* Main Content Styles */
.main-container {
    flex: 1;
    display: flex;
    flex-direction: column;
    position: relative;
    background-color: #ffffff;
    min-width: 0;
  }

/* Header */
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

.header-title-container {
  display: flex;
  align-items: center;
  gap: 8px;
  max-width: 240px; /* 进一步缩短最大宽度 */
  min-width: 0; /* 允许容器缩小 */
  padding: 6px 14px 6px 10px;
  background-color: #ffffff;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
  transition: all 0.2s ease;
  cursor: default;
  -webkit-app-region: no-drag;
}

.header-title-container:hover {
  background-color: #f9fafb;
  border-color: #d1d5db;
}

.header-title-icon {
  font-size: 16px;
  color: #6b7280;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.header-title-content {
  display: flex;
  align-items: center;
  overflow: hidden;
  flex: 1; /* 撑满剩余空间，超出部分隐藏 */
  min-width: 0; /* 必须加这个，否则 flex 内部的子元素无法缩小 */
}

.header-group-name {
  font-size: 13px;
  color: #6b7280;
  white-space: nowrap;
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

.header-separator {
  margin: 0 6px;
  color: #d1d5db;
  flex-shrink: 0;
}

.header-title {
  font-size: 14px;
  font-weight: 500;
  color: #111827;
  letter-spacing: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis; /* 文字过多直接省略号 */
  display: block; /* 配合 text-overflow 使用 */
}

.header-actions {
  flex-shrink: 0;
  margin-left: 16px;
  display: flex;
  align-items: center;
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
  -webkit-app-region: no-drag;
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
  -webkit-app-region: no-drag;
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

/* Content Area */
.content-area {
  flex: 1;
  overflow-y: auto;
  padding: 84px 0 160px 0; /* 恢复底部留白，避免遮挡 */
  display: flex;
  flex-direction: column;
  align-items: center;
  scroll-behavior: auto; /* 禁用全局平滑滚动，避免加载时闪烁或滑动 */
}

/* Initial Loading Skeleton */
.initial-tasks-loading {
  width: 100%;
  display: flex;
  flex-direction: column;
  padding: 0;
  gap: 32px; /* 与普通 message 的 gap 一致 */
  box-sizing: border-box;
}

.skeleton-group {
  display: flex;
  flex-direction: column;
  gap: 32px;
  width: 100%;
}

.skeleton-message {
  animation: pulse 1.5s infinite ease-in-out;
}

.skeleton-avatar-circle {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: #e5e7eb;
  flex-shrink: 0;
}

.skeleton-bg {
  background-color: #e5e7eb !important;
  color: transparent !important;
  border: none !important;
  box-shadow: none !important;
}

.user-text-skeleton { width: 40%; max-width: 300px; height: 24px; }
.user-text-skeleton.medium { width: 60%; max-width: 450px; }
.user-text-skeleton.short { width: 25%; max-width: 150px; }

.image-skeleton {
  width: 100%;
}
.image-skeleton.small { height: 200px; max-width: 200px; }
.image-skeleton.medium { height: 280px; max-width: 280px; }
.image-skeleton.large { height: 350px; max-width: 350px; }

@keyframes pulse {
  0% { opacity: 0.5; }
  50% { opacity: 0.8; }
  100% { opacity: 0.5; }
}

/* Welcome Screen */
.welcome-screen {
  max-width: 600px;
  width: 100%;
  margin-top: 10vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 0 20px;
  box-sizing: border-box;
}

.welcome-logo {
  margin-bottom: 24px;
  width: 80px;
  height: 80px;
  background-color: #f3f4f6;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.welcome-title {
  font-size: 32px;
  font-weight: 700;
  color: #111827;
  margin-bottom: 16px;
  margin-top: 0;
}

.welcome-desc {
  font-size: 16px;
  color: #6b7280;
  margin-bottom: 40px;
  line-height: 1.5;
}

.suggestion-cards {
  display: flex;
  gap: 16px;
  width: 100%;
  justify-content: center;
  flex-wrap: wrap;
}

.suggestion-card {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 20px;
  background-color: #f9fafb;
  border: 1px solid #e5e7eb;
  border-radius: 99px;
  font-size: 14px;
  color: #374151;
  cursor: pointer;
  transition: all 0.2s;
}

.suggestion-card:hover {
  background-color: #f3f4f6;
  border-color: #d1d5db;
}

/* Inspirations Panel */
.inspirations-panel {
  display: flex;
  flex-direction: column;
}

.inspirations-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid #e5e7eb;
  font-weight: 600;
  color: #374151;
}

.inspirations-content {
  padding: 8px;
  max-height: 300px;
  overflow-y: auto;
}

.inspiration-item {
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
  color: #4b5563;
  font-size: 13px;
  line-height: 1.5;
  margin-bottom: 6px;
}

.explore-area {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow-y: auto;
  background-color: #f9fafb;
}

/* Explore Dialog Styles */
.explore-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}
.explore-header {
  text-align: center;
  margin-bottom: 40px;
  padding-top: 20px;
}
.explore-header h2 {
  font-size: 32px;
  color: #111827;
  margin-bottom: 12px;
}
.explore-header p {
  font-size: 16px;
  color: #6b7280;
  margin-bottom: 24px;
}
.explore-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 24px;
}
.explore-card {
  background: #fff;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0,0,0,0.05);
  cursor: pointer;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  display: flex;
  flex-direction: column;
}
.explore-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 30px rgba(0,0,0,0.1);
}
.explore-image {
  width: 100%;
  aspect-ratio: 1 / 1;
  position: relative;
  background-color: #f3f4f6;
  overflow: hidden;
}
.explore-image .el-image {
  width: 100%;
  height: 100%;
}
.explore-no-image, .image-placeholder, .image-error {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  font-size: 32px;
}
.explore-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0,0,0,0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}
.explore-card:hover .explore-overlay {
  opacity: 1;
}
.explore-content {
  padding: 16px;
  font-size: 14px;
  color: #374151;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.inspiration-item:last-child {
  margin-bottom: 0;
}

.inspiration-item:hover {
  background-color: #f3f4f6;
  color: #111827;
}

.inspirations-popover {
  padding: 0 !important;
}

/* Result Screen (Chat style) */
.result-screen {
  width: 100%;
  max-width: 1000px; /* 增加聊天内容区的最大宽度 */
  padding: 0 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  gap: 32px;
  padding-bottom: 40px; /* 恢复适当空隙 */
}

.load-more-chat {
  display: flex;
  justify-content: center;
  margin: 10px 0 20px 0;
}

.load-more-inner {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 8px 20px;
  background-color: #f3f4f6;
  color: #6b7280;
  font-size: 13px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s ease;
  user-select: none;
}

.load-more-inner:hover {
  background-color: #e5e7eb;
  color: #3b82f6;
  transform: translateY(-1px);
}

.load-more-inner .mr-1 {
  margin-right: 2px;
  font-size: 14px;
}

.chat-time-divider {
  display: flex;
  justify-content: center;
  margin: 16px 0;
  width: 100%;
}

.chat-time-divider span {
  display: inline-block;
  padding: 4px 12px;
  background-color: rgba(0, 0, 0, 0.05);
  color: #999;
  font-size: 12px;
  border-radius: 12px;
  user-select: none;
}

.message {
  display: flex;
  width: 100%;
  margin-bottom: 0;
  animation: fadeIn 0.3s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

.message-inner {
  display: flex;
  gap: 16px;
  max-width: 95%; /* 增加消息气泡的最大宽度，减少空白 */
}

.user-message {
  justify-content: flex-end;
}

.user-message .message-inner {
  flex-direction: row;
}

.ai-message .message-inner {
  flex-direction: row;
  width: 100%; /* Make AI message inner take full width */
}

.message-content {
  padding: 16px 20px;
  border-radius: 20px;
  font-size: 15px;
  line-height: 1.6;
  word-break: break-word;
}

.user-content {
  background: linear-gradient(135deg, #ffffff 0%, #f3f4f6 100%);
  color: #1f2937;
  border-top-right-radius: 4px;
  border: 1px solid #e5e7eb;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
  cursor: default; /* Changed cursor to default since it has hover interaction */
}

.user-reference-image {
  margin-bottom: 12px;
  display: flex;
  flex-wrap: wrap;
  padding: 10px 10px 10px 0;
}

.reference-group {
  position: relative;
  transition: transform 0.2s, z-index 0s;
  cursor: pointer;
  margin-right: -15px; /* Overlap effect */
}

.reference-group:last-child {
  margin-right: 0;
}

.reference-group:hover {
  transform: translateY(-10px) scale(1.05);
  z-index: 99 !important;
}

.reference-item {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  border: 2px solid #ffffff;
  transition: all 0.3s;
  background: #ffffff;
}

.ref-img-display {
  width: 56px;
  height: 74px;
  display: block;
}

.message-meta {
  font-size: 12px;
  margin-top: 6px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  transition: opacity 0.2s ease;
  min-height: 20px; /* 保证悬停区域存在 */
}

.meta-date {
  color: #9ca3af;
  margin-right: 12px;
}

.meta-action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #9ca3af;
  font-size: 12px;
  cursor: pointer;
  transition: color 0.2s ease, opacity 0.2s ease;
  margin-right: 12px;
}

.meta-action-btn:hover {
  color: #3b82f6;
}

.meta-action-btn .el-icon {
  font-size: 14px;
}

.message-meta .meta-date,
.message-meta .meta-action-btn:not(.always-show),
.message-meta .simple-detail-btn {
  opacity: 0;
  transition: opacity 0.2s ease;
}

.message-meta .meta-action-btn.always-show {
  opacity: 1;
}

/* 扩大触发范围，只要在整条 message 里悬停就显示 */
.message:hover .message-meta .meta-date,
.message:hover .message-meta .meta-action-btn,
.message:hover .message-meta .simple-detail-btn,
.message-meta:hover .meta-date,
.message-meta:hover .meta-action-btn,
.message-meta:hover .simple-detail-btn,
.message-meta.is-active .meta-date,
.message-meta.is-active .meta-action-btn,
.message-meta.is-active .simple-detail-btn,
.el-popover[aria-hidden="false"] ~ * .message-meta .meta-date,
.el-popover[aria-hidden="false"] ~ * .message-meta .meta-action-btn,
.el-popover[aria-hidden="false"] ~ * .message-meta .simple-detail-btn,
.message-meta:focus-within .meta-date,
.message-meta:focus-within .meta-action-btn,
.message-meta:focus-within .simple-detail-btn {
  opacity: 1;
}

.simple-detail-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  color: #9ca3af;
  font-size: 12px;
  cursor: pointer;
  transition: color 0.2s ease;
}

.simple-detail-btn:hover {
  color: #3b82f6;
}

.simple-detail-btn .el-icon {
  font-size: 14px;
}

/* 详情弹窗样式 */
:deep(.custom-detail-popover) {
  padding: 0 !important;
  border-radius: 12px !important;
  overflow: hidden;
  border: 1px solid #e5e7eb !important;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1) !important;
}

.task-detail-content {
  display: flex;
  flex-direction: column;
}

.detail-header {
  padding: 12px 16px;
  background: #f8fafc;
  border-bottom: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  color: #334155;
  font-size: 14px;
}

.detail-header .el-icon {
  font-size: 16px;
  color: #3b82f6;
}

.detail-body {
  padding: 12px 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 13px;
}

.detail-row .label {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #64748b;
}

.detail-row .label .el-icon {
  font-size: 14px;
  opacity: 0.7;
}

.detail-row .value {
  color: #1e293b;
  font-weight: 500;
}

.highlight-row {
  margin-top: 4px;
  padding-top: 10px;
  border-top: 1px dashed #e2e8f0;
}

.highlight-value {
  color: #f59e0b !important;
  font-family: monospace;
  font-size: 14px;
}

.ai-avatar {
  width: 36px;
  height: 36px;
  background-color: #f3f4f6;
  color: #111827;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.ai-content {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.03);
  padding: 16px;
  border-top-left-radius: 4px;
  margin-top: 4px;
  width: 100%; /* Ensure content takes full width */
}

.generating-placeholder {
  position: relative;
  width: 100%;
  height: 100%;
  background-color: #f8fafc;
  border-radius: 16px;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  box-shadow: inset 0 0 20px rgba(0,0,0,0.02);
  min-height: 120px; /* Ensure there is a minimum height even if aspect ratio is squashed */
}

.placeholder-bg {
  position: absolute;
  inset: 0;
  background-image: linear-gradient(rgba(59, 130, 246, 0.05) 1px, transparent 1px),
                    linear-gradient(90deg, rgba(59, 130, 246, 0.05) 1px, transparent 1px);
  background-size: 20px 20px;
  background-position: center center;
  opacity: 0.8;
}

.placeholder-scan {
  position: absolute;
  inset: 0;
  background: linear-gradient(to bottom, transparent, rgba(59, 130, 246, 0.1) 50%, transparent);
  transform: translateY(-100%);
  animation: scan-effect 2.5s infinite ease-in-out;
}

@keyframes scan-effect {
  0% { transform: translateY(-100%); }
  100% { transform: translateY(100%); }
}

.placeholder-content {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  z-index: 2;
}

.generating-spinner {
  font-size: 28px;
  color: #3b82f6;
  filter: drop-shadow(0 0 8px rgba(59, 130, 246, 0.4));
}

.generating-text {
  font-size: 13px;
  font-weight: 500;
  color: #3b82f6;
  letter-spacing: 1px;
  animation: pulse-opacity 2s infinite ease-in-out;
}

.generating-timer {
  font-size: 12px;
  color: #64748b;
  font-family: monospace;
  margin-top: 4px;
}

@keyframes pulse-opacity {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}

.image-result-grid {
  width: 100%;
}

.image-grid-container {
  display: grid;
  gap: 12px;
  width: 100%;
}

.grid-count-1 {
  grid-template-columns: 1fr;
  max-width: 400px;
}

.grid-count-2, .grid-count-4 {
  grid-template-columns: repeat(2, 1fr);
  max-width: 560px;
}

.grid-count-3, .grid-count-5, .grid-count-6, .grid-count-7, .grid-count-8, .grid-count-9 {
  grid-template-columns: repeat(3, 1fr);
  max-width: 720px;
}

.image-grid-container:not(.grid-count-1):not(.grid-count-2):not(.grid-count-3):not(.grid-count-4):not(.grid-count-5):not(.grid-count-6):not(.grid-count-7):not(.grid-count-8):not(.grid-count-9) {
  grid-template-columns: repeat(4, 1fr);
  max-width: 860px;
}

.grid-item {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 100%;
  transition: transform 0.2s ease;
}

.grid-item:hover {
  transform: translateY(-2px);
}

.generated-image {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
  min-height: 120px;
  background: #f9fafb;
  border: 1px solid #f3f4f6;
}

.generated-image :deep(.el-image__error) {
  height: 100%;
  width: 100%;
  display: flex;
}

.image-actions {
  display: flex;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s ease, transform 0.2s ease;
  position: absolute;
  bottom: 8px;
  right: 8px;
  background: rgba(255, 255, 255, 0.95);
  padding: 0 4px;
  border-radius: 8px;
  backdrop-filter: blur(4px);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  transform: translateY(4px);
}

.grid-item:hover .image-actions {
  opacity: 1;
  transform: translateY(0);
}

.image-slot {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  min-height: 150px;
  background: #f9fafb;
  color: #9ca3af;
  font-size: 30px;
  border-radius: 12px;
  border: 1px solid #f3f4f6;
  box-sizing: border-box;
}

/* Input Area */
.input-area-wrapper {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  width: 100%;
  background: transparent;
  padding: 10px 0 calc(32px + env(safe-area-inset-bottom)) 0; /* 整体调高，避免沉底 */
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-shrink: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: bottom center;
  z-index: 100;
  pointer-events: none;
}

.input-area-wrapper .input-container,
.input-area-wrapper .input-tools-container {
  pointer-events: auto;
}

.input-area-wrapper.is-shrunk {
  padding-bottom: calc(32px + env(safe-area-inset-bottom)) !important; /* 增加悬浮高度 */
}

.input-area-wrapper.is-shrunk::before {
  opacity: 0;
}

.input-area-wrapper.is-shrunk .input-tools-container.mobile-only,
.input-area-wrapper.is-shrunk .input-tools.desktop-only {
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
  
  .reference-upload-item.initial-upload-item {
    width: 40px;
    height: 48px;
    margin-right: 12px; /* 适当增加初始状态下旋转卡片与文字的间距，防止遮挡 */
  }
  
  .preview-list-inline {
    width: 48px; /* 恢复为电脑端的宽度 */
    height: 72px; /* 恢复为电脑端的高度 */
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

.input-tools-container {
  width: calc(100% - 40px);
  max-width: 760px;
  margin-bottom: 8px;
  display: flex;
  justify-content: flex-start;
}

.input-container {
  width: calc(100% - 40px);
  max-width: 760px;
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
}

.model-settings-btn {
  width: 100%;
  min-width: 100px;
  max-width: 200px;
  justify-content: space-between;
}

.model-btn-content {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  min-width: 0;
}

.cpu-icon {
  flex-shrink: 0;
}

.model-name-text {
  display: block !important;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
  text-align: left;
  line-height: 1.5;
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

.input-tools-container.mobile-only,
.input-tools.desktop-only {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: bottom center;
}

.input-tools {
  display: flex;
  gap: 8px;
  padding: 0;
  align-items: center;
  flex-wrap: nowrap;
  overflow-x: auto;
  overflow-y: hidden;
  scrollbar-width: none;
}
.input-tools::-webkit-scrollbar {
  display: none;
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
  max-width: 800px;
  background: transparent;
  border: none;
  border-radius: 20px;
  padding: 8px 4px;
  gap: 12px;
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
}

:deep(.chat-input .el-textarea__inner) {
  padding-left: 8px !important;
  box-shadow: none !important;
  background-color: transparent !important;
  padding-top: 0px !important;
  padding-bottom: 0px !important;
  font-size: 15px;
  line-height: 24px !important; /* 修改为固定的行高，使单行文本垂直居中 */
  height: 24px !important; /* 确保基础高度匹配行高 */
  min-height: 24px !important;
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
@media (min-width: 769px) {
  .header-actions {
    display: none;
  }
  .mobile-menu-btn {
    display: none;
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
    max-width: 95%;
    width: 100%;
    padding: 0 4px;
    box-sizing: border-box;
    margin-bottom: 0;
  }

  .input-container {
    max-width: 95%;
    width: 100%;
    padding: 12px;
    border-radius: 16px;
  }

  .input-area-wrapper {
    padding: 10px 0 10px 0;
  }

  .message-inner {
    max-width: 95%;
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
    overflow-y: hidden;
    scrollbar-width: none !important; /* Firefox 强制隐藏滚动条 */
    -ms-overflow-style: none; /* IE and Edge 强制隐藏滚动条 */
    -webkit-overflow-scrolling: touch;
    margin-bottom: 8px;
  }
  
  /* Webkit 浏览器 (Chrome, Safari, iOS 等) 强制隐藏滚动条 */
  .input-tools::-webkit-scrollbar {
    display: none !important;
    width: 0 !important;
    height: 0 !important;
    background: transparent !important;
  }

  .model-dropdown-wrapper {
    max-width: 100%;
    flex: 0 0 auto; /* 改为固定大小，不允许压缩，以支持横向滑动 */
    min-width: 0;
  }
  
  .combined-settings-btn {
    padding: 4px 12px;
    height: 32px;
    font-size: 13px;
    max-width: 100%;
    width: auto;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    box-sizing: border-box;
    flex: 0 0 auto;
    border-radius: 16px;
    background-color: #ffffff; /* 增加背景色 */
    border: 1px solid #e5e7eb; /* 增加边框 */
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05); /* 增加轻微阴影提升立体感 */
  }

  .model-settings-btn {
    justify-content: flex-start; /* 模型名字靠左，箭头靠右（如果有足够空间） */
  }

  .model-btn-content {
    flex: 0 1 auto;
    min-width: 0; /* 允许文本截断 */
  }

  .model-name-text {
    flex: 0 1 auto;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    max-width: 160px; /* 在手机端给个合理的最大宽度 */
    line-height: 1.5;
  }

  .combined-settings-btn .divider {
    margin: 0 6px;
  }

  .upload-btn-wrapper {
    flex: 0 0 auto;
  }

  .ref-btn {
    height: 32px;
    font-size: 13px;
    padding: 0 12px;
    border-radius: 16px;
  }
  
  /* 手机端图片网格排版优化：智能分配列数 */
  .image-grid-container:not(.grid-count-1) {
    grid-template-columns: repeat(2, 1fr) !important;
    gap: 8px !important;
  }
  
  /* 只有 3, 6, 9 张图时，使用 3 列以保证完美对齐，其余多图统一使用 2 列以保证图片够大 */
  .image-grid-container.grid-count-3,
  .image-grid-container.grid-count-6,
  .image-grid-container.grid-count-9 {
    grid-template-columns: repeat(3, 1fr) !important;
  }
  
  .generated-image {
    min-height: 80px;
    border-radius: 8px;
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
  padding: 4px 14px;
  height: 34px;
  font-size: 13px;
  color: #3f3f46;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: none;
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
  line-height: normal;
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

@media (max-width: 480px) {
  .premium-redeem-container {
    padding: 32px 24px 24px;
  }
  
  .redeem-icon-wrapper {
    width: 56px;
    height: 56px;
    margin-bottom: 16px;
  }
  
  .redeem-icon-wrapper svg {
    width: 28px;
    height: 28px;
  }
  
  .redeem-title {
    font-size: 20px;
  }
  
  .redeem-desc {
    font-size: 13px;
    margin-bottom: 24px;
  }
  
  .custom-redeem-input {
    height: 48px;
    font-size: 15px;
  }
  
  .redeem-body {
    margin-bottom: 24px;
  }
  
  .redeem-footer {
    gap: 10px;
  }
  
  .redeem-action-btn {
    height: 42px;
    font-size: 14px;
  }
}

/* 移动到项目组弹窗样式 - 现代化质感UI */
.modern-dialog-header {
  padding: 24px 24px 16px;
}

.modern-dialog-title {
  margin: 0 0 6px 0;
  font-size: 18px;
  font-weight: 600;
  color: #111827;
}

.modern-dialog-subtitle {
  margin: 0;
  font-size: 14px;
  color: #6b7280;
}

.modern-group-list-wrapper {
  padding: 0 24px;
}

.modern-group-list {
  max-height: 280px;
  overflow-y: auto;
  margin: 0 -4px;
  padding: 4px;
}

.modern-group-list::-webkit-scrollbar {
  width: 4px;
}
.modern-group-list::-webkit-scrollbar-thumb {
  background: #e5e7eb;
  border-radius: 2px;
}

.modern-group-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 8px;
  background-color: #ffffff;
  border: 1px solid #f3f4f6;
}

.modern-group-item:last-child {
  margin-bottom: 0;
}

.modern-group-item:hover {
  border-color: #e5e7eb;
  background-color: #f9fafb;
}

.modern-group-item.is-active {
  background-color: #ffffff;
  border-color: #3b82f6;
  box-shadow: 0 0 0 1px #3b82f6;
}

.modern-group-item-left {
  display: flex;
  align-items: center;
  gap: 12px;
  overflow: hidden;
}

.modern-group-icon-wrapper {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background-color: #f3f4f6;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
}

.modern-group-item.is-active .modern-group-icon-wrapper {
  background-color: #eff6ff;
}

.modern-group-icon {
  font-size: 18px;
  color: #6b7280;
  transition: color 0.2s;
}

.modern-group-item.is-active .modern-group-icon {
  color: #3b82f6;
}

.modern-group-name {
  font-size: 15px;
  color: #374151;
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.modern-group-item.is-active .modern-group-name {
  color: #111827;
}

.modern-check-wrapper {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 1.5px solid #d1d5db;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  margin-right: 4px;
}

.modern-check-wrapper.is-checked {
  background-color: #3b82f6;
  border-color: #3b82f6;
}

.modern-check-wrapper .el-icon {
  color: white;
  font-size: 12px;
}

.modern-group-empty {
  padding: 20px 0;
}

.modern-dialog-footer {
  padding: 24px;
  display: flex;
  gap: 12px;
}

.modern-btn {
  flex: 1;
  height: 44px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modern-btn-cancel {
  background-color: #ffffff;
  border: 1px solid #d1d5db;
  color: #374151;
}

.modern-btn-cancel:hover {
  background-color: #f9fafb;
}

.modern-btn-confirm {
  background-color: #3b82f6;
  color: #ffffff;
}

.modern-btn-confirm:hover:not(:disabled) {
  background-color: #2563eb;
}

.modern-btn-confirm:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
}

@media (max-width: 768px) {
  .modern-dialog-header {
    padding: 20px 20px 12px;
  }
  .modern-group-list-wrapper {
    padding: 0 20px;
  }
  .modern-dialog-footer {
    padding: 20px;
  }
}
</style>

<style>
.el-dialog.modern-group-dialog {
  border-radius: 16px !important;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25) !important;
  padding: 0 !important;
  overflow: hidden;
}

.el-dialog.modern-group-dialog .el-dialog__header {
  display: none !important;
}
.el-dialog.modern-group-dialog .el-dialog__body {
  padding: 0 !important;
}

@media (max-width: 768px) {
  .el-dialog.modern-group-dialog {
    width: 90vw !important;
    max-width: 400px;
  }
}

.el-dialog.premium-redeem-dialog {
  border-radius: 24px !important;
  overflow: hidden !important;
  background: #ffffff !important;
  padding: 0 !important;
  box-shadow: 0 24px 48px rgba(0, 0, 0, 0.15) !important;
  max-width: 90vw !important;
}

@media (max-width: 480px) {
  .el-dialog.premium-redeem-dialog {
    width: 340px !important;
  }
}

.el-dialog.premium-redeem-dialog .el-dialog__header {
  display: none !important;
}
.el-dialog.premium-redeem-dialog .el-dialog__body {
  padding: 0 !important;
}
</style>
