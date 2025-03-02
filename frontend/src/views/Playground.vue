<script setup lang="ts">
import {Button} from '@/components/ui/button';
import {ScrollArea} from '@/components/ui/scroll-area';
import {Badge} from '@/components/ui/badge';
import {ResizableHandle, ResizablePanel, ResizablePanelGroup} from '@/components/ui/resizable';
import {Card} from '@/components/ui/card';
import {Play, RotateCcw} from "lucide-vue-next"
import VCodeBlock from '@wdns/vue-code-block';

import {useToast} from '@/components/ui/toast/use-toast'

import {computed, ref} from "vue";
import {useRouter} from "vue-router";
import {engines, main} from "../../wailsjs/go/models.ts";

import mime from "mime"
import {GeneratePatchedFile, ReadFile, SelectFile} from "../../wailsjs/go/main/App";
import Settings from "@/components/Settings.vue";
import LanguageIcon from "@/components/icons/LanguageIcon.vue";
import GeneratePatchedFileArgs = main.GeneratePatchedFileArgs;

interface ConfigurationProps {
  engines: engines.EngineMeta[],
}

const fileContent = ref("");

const {toast} = useToast()

const router = useRouter()

const props = defineProps<ConfigurationProps>()

const currentEngine = computed(() => {
  for (let i = 0; i < props.engines.length; i++) {
    if (props.engines[i].name.toLowerCase() === (router.currentRoute.value.params.engineName as string)?.toLowerCase()) {
      Object.entries(props.engines[i].args as engines.EngineArg[]).forEach(([_, arg]) => {
        arg.value = arg.defaultValue
      })
      return props.engines[i]
    }
  }
  return props.engines[0]
})

const originalDroppedFile = ref<string>("")
const currentFilePath = ref<string | null>(null)
const isDragging = ref(false)

const handleDragOver = (e: DragEvent) => {
  e.preventDefault()
  isDragging.value = true
}

const handleDragLeave = (e: DragEvent) => {
  e.preventDefault()
  isDragging.value = false
}

const handleDrop = async (e: DragEvent) => {
  e.preventDefault()
  isDragging.value = false

  const files = e.dataTransfer?.files
  if (!files?.length) return

  const file = files[0]
  currentFilePath.value = file.name

  try {
    fileContent.value = await readDroppedFile(file)
    originalDroppedFile.value = fileContent.value
  } catch (err) {
    fileContent.value = ""
    currentFilePath.value = ""
    toast({
      title: 'Error reading file' as string,
      description: err as string,
      variant: "destructive"
    })

    console.error('Error reading file:', err)
  }
}

const handleFileSelect = async () => {
  try {
    const res = await SelectFile()
    if (res) {
      currentFilePath.value = res
    } else {
      return
    }

    try {
      fileContent.value = await ReadFile(currentFilePath.value as string);
      originalDroppedFile.value = fileContent.value
    } catch (err) {
      fileContent.value = ""
      currentFilePath.value = ""
      toast({
        title: 'Error reading file' as string,
        description: err as string,
        variant: "destructive"
      })

      console.error('Error reading file:', err)
    }

  } catch (error) {
    console.error('Error selecting file:', error)
  }
}

function truncateFilePath(filePath: string, depth: number): string {
  const parts = filePath.split(/[/\\]/); // Split by both forward and backslashes
  const truncatedParts = parts.slice(-depth); // Get the last 3 elements
  return parts.length > depth ? `.../${truncatedParts.join("/")}` : filePath;
}

const applyEngine = () => {
  if (!fileContent.value) {
    toast({
      title: 'Failed to run' as string,
      description: 'No file loaded' as string,
    })
    return
  }

  const args = <GeneratePatchedFileArgs>{
    engine: currentEngine.value.name,
    content: fileContent.value,
    settings: currentEngine.value.args
  }

  GeneratePatchedFile(args).then((patched) => {
    fileContent.value = patched
  }).catch((err) => {
    toast({
      title: 'Error patching the file' as string,
      description: err as string,
      variant: "destructive"
    })
  })

}

const readDroppedFile = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()

    reader.onload = (event) => {
      const content = event.target?.result as string
      if (content.length > 500000) {
        reject(new Error("File size exceeds the 500KB limit"))
      }

      resolve(content)
    }

    reader.onerror = (error) => {
      reject(error)
    }

    reader.readAsText(file)
  })
}

const parseLangFromExt = (filepath: string) => {
  if (filepath.substring(filepath.lastIndexOf('.')) === '.go') {
    return 'golang';
  }else  if (filepath.substring(filepath.lastIndexOf('.')) === '.ps1') {
    return 'powershell';
  }

  const mimeType = mime.getType(filepath);
  // at least try 🙆‍♀️
  const extension = mime.getExtension(mimeType || '')?.replace('x-', '');

  return extension || 'plain';
}
</script>

<template>
  <div class="grid flex-1 gap-4 overflow-auto">
    <main class="grid flex-1 gap-4 overflow-auto p-4 pb-0 ">
      <ResizablePanelGroup direction="horizontal">
        <ResizablePanel :default-size="25" :min-size="10">
          <div class="grid flex-1 gap-4 overflow-auto pt-2 pb-4 pl-1 pr-6 lg:p-4 lg:pl-2 lg:pr-6 h-full">
            <ScrollArea class="flex-1">
              <div
                  class="mb-4 px-2"
              >
                <span
                    class="flex justify-between items-center"
                >
                  <span
                      class="text-3xl font-extrabold tracking-[-0.04em]"
                  >
                    {{ currentEngine.name }}
                  </span>
                  <span
                    class="h-full"
                  >
                      <LanguageIcon :lang="currentEngine.language"/>
                  </span>
                </span>
                <p
                    class="text-base/7 text-foreground/80"
                >
                  {{ currentEngine.longDescription }}
                </p>
              </div>
              <Settings
                  :current-engine="currentEngine"
                  :args="currentEngine?.args"
              />
            </ScrollArea>
            <div
                class="flex justify-between mt-auto flex-col lg:flex-row"
            >
              <Button
                  class="cursor-pointer w-full"
                  size="lg"
                  @click="applyEngine"
              >
                <Play/>
                Run
              </Button>
              <Button
                  class="cursor-pointer mt-4 lg:mt-0 lg:ml-4 bg-secondary text-secondary-foreground hover:bg-primary/10"
                  size="lg"
                  @click="fileContent = originalDroppedFile"
              >
                <RotateCcw/>
                Reset
              </Button>
            </div>
          </div>
        </ResizablePanel>
        <ResizableHandle with-handle/>
        <ResizablePanel :default-size="75" :min-size="20">
          <div :class="['relative flex z-50 h-full flex-col rounded-xl p-4 lg:col-span-2']"
               @dragover="(e) =>  handleDragOver(e)"
               @dragleave="(e) => handleDragLeave(e)"
               @drop="(e) => handleDrop(e)"
          >
            <Badge
                variant="outline"
                class="absolute right-2 top-2 z-40 bg-secondary font-medium"
            >
              {{ currentFilePath ? truncateFilePath(currentFilePath, 3) : 'Output' }}
            </Badge>
            <div
                v-if="fileContent.length > 0"
                :class="['flex-1 h-[75dvh] bg-[#282c34] rounded-xl', isDragging ? 'grayscale' : '']"
            >
              <ScrollArea class=" h-[calc(100vh_-_110px)] w-full rounded-xl"
              >
                <VCodeBlock
                    :class="['w-full h-full rounded-xl', isDragging ? 'border-red-600 grayscale':'border-blue-600']"
                    :code="fileContent"
                    persistent-copy-button
                    highlightjs
                    :lang="parseLangFromExt(currentFilePath as string)"
                    theme="atom-one-dark"
                />
              </ScrollArea>
            </div>
            <div v-else class="flex-1 z-10 bg-muted/50 rounded-xl">
              <div v-if="!currentFilePath"
                   class="h-full flex items-center justify-center mx-4">
                <div class="w-96">
                  <Card
                      @click="handleFileSelect"
                      :class="[
                        'p-6 border-2 border-dashed transition-colors duration-200 cursor-pointer',
                        isDragging ? 'border-primary bg-muted/10' : 'border-muted-foreground/25'
                      ]">
                    <div class="text-center space-y-2">
                      <div class="text-3xl font-semibold">
                        +
                      </div>
                      <p class="text-sm text-muted-foreground">
                        Drag and drop a file here to begin
                      </p>
                    </div>
                  </Card>
                </div>
              </div>
            </div>
          </div>
        </ResizablePanel>
      </ResizablePanelGroup>
    </main>
  </div>

</template>

<style scoped>

</style>