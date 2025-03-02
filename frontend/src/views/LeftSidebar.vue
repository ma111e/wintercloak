<script setup lang="ts">
import {Collapsible, CollapsibleContent, CollapsibleTrigger} from '@/components/ui/collapsible'
import {Tabs, TabsList, TabsTrigger} from '@/components/ui/tabs'
import {Icon} from '@iconify/vue'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarProvider,
} from '@/components/ui/sidebar'
import {ChevronRight, Cog, Home, Wrench, Paintbrush} from 'lucide-vue-next'
import {engines} from "../../wailsjs/go/models.ts";
import {useRouter} from "vue-router";

// import LanguageIcon from "@/components/icons/LanguageIcon.vue";

interface SidebarProps {
  engines: engines.EngineMeta[]
}

defineProps<SidebarProps>()
const router = useRouter()

const data = {
  navEngines: [
    {
      title: 'Engines',
      url: '',
      icon: Cog,
      isActive: true,
    }
  ],
  navGeneric: [
      {
      title: 'Settings',
      icon: Wrench,
      isActive: true,
      items: [
        {
          name: 'Theme',
          url: '/settings/theme',
          icon: Paintbrush,
        },
      ]
    }
  ],
}

</script>

<template>
  <SidebarProvider
      class="w-65 h-[calc(100vh_-_57px)] min-h-0"
  >
    <Sidebar
        class="pt-[50px] w-65"
    >
      <SidebarHeader>

      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup class="group-data-[collapsible=icon]:hidden">
          <SidebarMenu>
            <SidebarMenuItem>
              <SidebarMenuButton as-child>
                <router-link
                    :class="['', (router.currentRoute.value.path as string) === '/' ? 'bg-muted': '']"
                    to="/"
                >
                  <Home/>
                  <span>Home</span>
                </router-link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
          <SidebarMenu>
            <Collapsible
                v-for="item in data.navEngines"
                :key="item.title"
                as-child
                :default-open="item.isActive"
                class="group/collapsible"
            >
              <SidebarMenuItem>
                <CollapsibleTrigger as-child>
                  <SidebarMenuButton :tooltip="item.title">
                    <component :is="item.icon"/>
                    <span>{{ item.title }}</span>
                    <ChevronRight
                        class="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"/>
                  </SidebarMenuButton>
                </CollapsibleTrigger>
                <CollapsibleContent>
                  <SidebarMenuSub>
                    <SidebarMenuSubItem
                        v-for="subitem in $props.engines"
                        :key="subitem.name"
                    >
                      <SidebarMenuSubButton as-child>
                        <router-link
                            :class="['flex items-center text-muted-foreground hover:bg-muted hover:text-muted-foreground', (router.currentRoute.value.params.engineName as string)?.toLowerCase() === subitem.name.toLowerCase() ? 'bg-muted': '']"
                            :to="{name: 'engines', params:{ engineName: `${subitem.name.toLowerCase()}`}, query: {...$route.query}}"
                        >
                          <Icon :icon="subitem.iconifyIcon"/>
                          <!--                          <LanguageIcon :lang="eng.language"/>-->
                          <span>{{ subitem.name }}</span>
                        </router-link>
                      </SidebarMenuSubButton>
                    </SidebarMenuSubItem>
                  </SidebarMenuSub>
                </CollapsibleContent>
              </SidebarMenuItem>
            </Collapsible>
          </SidebarMenu>
          <SidebarMenu>
            <Collapsible
                v-for="item in data.navGeneric"
                :key="item.title"
                as-child
                :default-open="item.isActive"
                class="group/collapsible"
            >
              <SidebarMenuItem>
                <CollapsibleTrigger as-child>
                  <SidebarMenuButton :tooltip="item.title">
                    <component :is="item.icon"/>
                    <span>{{ item.title }}</span>
                    <ChevronRight
                        class="ml-auto transition-transform duration-200 group-data-[state=open]/collapsible:rotate-90"/>
                  </SidebarMenuButton>
                </CollapsibleTrigger>
                <CollapsibleContent>
                  <SidebarMenuSub>
                    <SidebarMenuSubItem
                        v-for="subitem in item.items"
                        :key="subitem.name"
                    >
                      <SidebarMenuSubButton as-child>
                        <router-link
                            :class="['flex items-center text-muted-foreground hover:bg-muted hover:text-muted-foreground', (router.currentRoute.value.path as string)?.toLowerCase() === subitem.url.toLowerCase() ? 'bg-muted': '']"
                            :to="subitem.url"
                        >
                          <!--                            :to="`/engines/${eng.name.toLowerCase()}`"-->

                          <component :is="subitem.icon"/>
                          <!--                          <LanguageIcon :lang="eng.language"/>-->
                          <span>{{ subitem.name }}</span>
                        </router-link>
                      </SidebarMenuSubButton>
                    </SidebarMenuSubItem>
                  </SidebarMenuSub>
                </CollapsibleContent>
              </SidebarMenuItem>
            </Collapsible>
          </SidebarMenu>
        </SidebarGroup>
        <!--        <SidebarGroup class="group-data-[collapsible=icon]:hidden">-->
        <!--          <SidebarGroupLabel>Resources</SidebarGroupLabel>-->
        <!--          <SidebarMenu>-->
        <!--            <SidebarMenuItem-->
        <!--                v-for="item in data.resources"-->
        <!--                :key="item.name"-->
        <!--            >-->
        <!--              <SidebarMenuButton as-child>-->
        <!--                <router-link-->
        <!--                    :to="item.url"-->
        <!--                >-->
        <!--                  <component :is="item.icon"/>-->
        <!--                  <span>{{ item.name }}</span>-->
        <!--                </router-link>-->
        <!--                &lt;!&ndash;                <a :href="item.url">&ndash;&gt;-->
        <!--                &lt;!&ndash;                  <component :is="item.icon"/>&ndash;&gt;-->
        <!--                &lt;!&ndash;                  <span>{{ item.name }}</span>&ndash;&gt;-->
        <!--                &lt;!&ndash;                </a>&ndash;&gt;-->
        <!--              </SidebarMenuButton>-->
        <!--            </SidebarMenuItem>-->
        <!--          </SidebarMenu>-->
        <!--        </SidebarGroup>-->
      </SidebarContent>
      <SidebarFooter>
        <Tabs
            v-if="$router.currentRoute.value.fullPath.startsWith('/engines/')"
            :default-value="$route.query.mode as string || 'playground'"
            class="w-full"
        >
          <TabsList>
            <!--                @click="$emit('setPlayground', false)"-->
            <TabsTrigger
                class="hover:bg-primary/10 cursor-pointer"
                value="preprocessor"
                @click="$router.replace({ query: { ...$route.query, mode: 'preprocessor' } })"
            >
              Preprocessor
            </TabsTrigger>

            <TabsTrigger
                class="hover:bg-primary/10 cursor-pointer"
                value="playground"
                @click="$router.replace({ query: { ...$route.query, mode: 'playground' } })"
            >
              Playground
            </TabsTrigger>
          </TabsList>
        </Tabs>
      </SidebarFooter>
    </Sidebar>
  </SidebarProvider>
</template>