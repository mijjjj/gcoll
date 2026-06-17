import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import DevicesView from '../views/DevicesView.vue'
import LogsView from '../views/LogsView.vue'
import PipelineView from '../views/PipelineView.vue'
import PluginsView from '../views/PluginsView.vue'
import PointCacheView from '../views/PointCacheView.vue'
import TargetsView from '../views/TargetsView.vue'
import TasksView from '../views/TasksView.vue'

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/dashboard',
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: DashboardView,
    },
    {
      path: '/devices',
      name: 'devices',
      component: DevicesView,
    },
    {
      path: '/tasks',
      name: 'tasks',
      component: TasksView,
    },
    {
      path: '/point-cache',
      name: 'point-cache',
      component: PointCacheView,
    },
    {
      path: '/pipeline',
      name: 'pipeline',
      component: PipelineView,
    },
    {
      path: '/targets',
      name: 'targets',
      component: TargetsView,
    },
    {
      path: '/plugins',
      name: 'plugins',
      component: PluginsView,
    },
    {
      path: '/logs',
      name: 'logs',
      component: LogsView,
    },
  ],
})
