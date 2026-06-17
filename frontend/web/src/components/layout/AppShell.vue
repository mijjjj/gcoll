<script setup lang="ts">
import { computed, h, ref, watchEffect } from 'vue'
import {
  darkTheme,
  lightTheme,
  NButton,
  NConfigProvider,
  NDialogProvider,
  NGlobalStyle,
  NLayout,
  NLayoutContent,
  NLayoutHeader,
  NLayoutSider,
  NMessageProvider,
  NMenu,
  NTooltip,
  type GlobalThemeOverrides,
  type MenuOption,
} from 'naive-ui'
import {
  Activity,
  Bell,
  ClipboardList,
  Database,
  FileText,
  Globe2,
  Info,
  Languages,
  ListChecks,
  Minus,
  Moon,
  PanelLeftClose,
  Plug,
  Puzzle,
  Settings,
  Square,
  Sun,
  Wifi,
  X,
} from '@lucide/vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import HttpErrorBridge from '../common/HttpErrorBridge.vue'
import { useLocaleStore } from '../../i18n'

const route = useRoute()
const router = useRouter()
const localeStore = useLocaleStore()
const storedTheme = localStorage.getItem('gcoll-theme')
const darkMode = ref(storedTheme === 'dark')
const t = computed(() => localeStore.t)

const theme = computed(() => (darkMode.value ? darkTheme : lightTheme))
const themeOverrides = computed<GlobalThemeOverrides>(() => ({
  common: {
    primaryColor: darkMode.value ? '#38BDF8' : '#0E7CF2',
    primaryColorHover: darkMode.value ? '#7DD3FC' : '#0B66D8',
    primaryColorPressed: darkMode.value ? '#0EA5E9' : '#0757C2',
    borderRadius: '6px',
    fontSize: '13px',
    heightSmall: '28px',
    heightMedium: '32px',
  },
  Button: {
    borderRadiusMedium: '6px',
  },
  DataTable: {
    thColor: darkMode.value ? '#172033' : '#F7FAFD',
    tdColor: darkMode.value ? '#111827' : '#FFFFFF',
    borderColor: darkMode.value ? '#243244' : '#DCE7F2',
  },
  Menu: {
    itemHeight: '40px',
    itemIconSize: '17px',
    itemTextColorActive: darkMode.value ? '#E5EEF8' : '#0E7CF2',
    itemIconColorActive: darkMode.value ? '#38BDF8' : '#0E7CF2',
    itemColorActive: darkMode.value ? '#083344' : '#EAF4FF',
  },
}))

const renderIcon = (icon: typeof ListChecks) => () => h(icon, { size: 17, strokeWidth: 1.8 })

const menuOptions = computed<MenuOption[]>(() => [
  { label: t.value('nav.dashboard'), key: 'dashboard', icon: renderIcon(Activity) },
  { label: t.value('nav.devices'), key: 'devices', icon: renderIcon(ListChecks) },
  { label: t.value('nav.tasks'), key: 'tasks', icon: renderIcon(ClipboardList) },
  { label: t.value('nav.data'), key: 'point-cache', icon: renderIcon(FileText) },
  { label: t.value('nav.pipeline'), key: 'pipeline', icon: renderIcon(Settings) },
  { label: t.value('nav.targets'), key: 'targets', icon: renderIcon(Globe2) },
  { label: t.value('nav.plugins'), key: 'plugins', icon: renderIcon(Puzzle) },
  { label: t.value('nav.logs'), key: 'logs', icon: renderIcon(Plug) },
  { label: t.value('nav.alerts'), key: 'alerts', icon: renderIcon(Bell), disabled: true },
  { label: t.value('nav.about'), key: 'about', icon: renderIcon(Info), disabled: true },
])

const activeKey = computed(() => route.name?.toString() ?? 'dashboard')

function handleMenuUpdate(key: string) {
  router.push({ name: key })
}

watchEffect(() => {
  localStorage.setItem('gcoll-theme', darkMode.value ? 'dark' : 'light')
})
</script>

<template>
  <NConfigProvider :theme="theme" :theme-overrides="themeOverrides">
    <NGlobalStyle />
    <NMessageProvider>
      <NDialogProvider>
        <HttpErrorBridge />
        <NLayout class="app-shell" :class="{ 'is-dark': darkMode }" has-sider>
          <NLayoutSider bordered :width="162" :native-scrollbar="false" class="app-sider">
            <div class="brand">
              <div class="brand-mark">g</div>
              <strong>gcoll</strong>
            </div>

            <NMenu
              :value="activeKey"
              :options="menuOptions"
              :root-indent="20"
              :indent="12"
              class="side-menu"
              @update:value="handleMenuUpdate"
            />

            <div class="sider-bottom">
              <NButton quaternary size="small" class="collapse-btn">
                <template #icon>
                  <PanelLeftClose :size="16" />
                </template>
                {{ t('sidebar.collapse') }}
              </NButton>
              <span>v1.2.0</span>
            </div>
          </NLayoutSider>

          <NLayout class="app-main">
            <NLayoutHeader class="topbar">
              <div class="status-card status-card--db">
                <div class="status-icon"><Database :size="22" /></div>
                <div class="status-copy">
                  <strong>{{ t('status.sqlite') }} <span class="status-dot-text">{{ t('common.running') }}</span></strong>
                  <small>~/.gcoll/data/gcoll.db</small>
                </div>
              </div>

              <div class="status-card status-card--api">
                <div class="status-icon"><Globe2 :size="22" /></div>
                <div class="status-copy">
                  <strong>{{ t('status.httpApi') }} <span class="status-dot-text">{{ t('common.running') }}</span></strong>
                  <small>http://127.0.0.1:4120</small>
                </div>
                <div class="api-key-chip">APIKey</div>
                <div class="api-secret">sk_live_••••••••••••••</div>
              </div>

              <div class="status-card status-card--plugin">
                <div class="status-icon"><Puzzle :size="22" /></div>
                <div class="status-copy">
                  <strong>{{ t('status.pluginProcess') }}</strong>
                  <small class="success-text">{{ t('status.pluginCount') }}</small>
                </div>
              </div>

              <div class="status-card status-card--network">
                <div class="status-icon"><Wifi :size="22" /></div>
                <div class="status-copy">
                  <strong>{{ t('status.network') }}</strong>
                  <small class="warning-text">{{ t('common.offlineMode') }}</small>
                </div>
              </div>

              <div class="theme-switch">
                <NButton
                  size="small"
                  secondary
                  :aria-label="darkMode ? t('common.theme.switchToLight') : t('common.theme.switchToDark')"
                  @click="darkMode = !darkMode"
                >
                  <template #icon>
                    <component :is="darkMode ? Sun : Moon" :size="14" />
                  </template>
                  {{ darkMode ? t('common.theme.dark') : t('common.theme.light') }}
                </NButton>
                <NButton size="small" secondary :aria-label="t('common.language.switch')" @click="localeStore.toggleLanguage()">
                  <template #icon>
                    <Languages :size="14" />
                  </template>
                  {{ t('common.language.current') }}
                </NButton>
              </div>

              <div class="window-actions">
                <NTooltip trigger="hover">
                  <template #trigger>
                    <NButton quaternary circle size="small" :aria-label="t('common.minimize')">
                      <Minus :size="15" />
                    </NButton>
                  </template>
                  {{ t('common.minimize') }}
                </NTooltip>
                <NTooltip trigger="hover">
                  <template #trigger>
                    <NButton quaternary circle size="small" :aria-label="t('common.maximize')">
                      <Square :size="13" />
                    </NButton>
                  </template>
                  {{ t('common.maximize') }}
                </NTooltip>
                <NTooltip trigger="hover">
                  <template #trigger>
                    <NButton quaternary circle size="small" :aria-label="t('common.close')">
                      <X :size="16" />
                    </NButton>
                  </template>
                  {{ t('common.close') }}
                </NTooltip>
              </div>
            </NLayoutHeader>

            <NLayoutContent class="content">
              <RouterView />
            </NLayoutContent>
          </NLayout>
        </NLayout>
      </NDialogProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>
