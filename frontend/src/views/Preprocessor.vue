<script setup lang="ts">
import {Label} from "@/components/ui/label";
import {Separator} from "@/components/ui/separator";
import {computed, ref} from "vue";
import {useRouter} from "vue-router";
import {Play} from "lucide-vue-next"
import {Button} from "@/components/ui/button";
import {ScrollArea} from "@/components/ui/scroll-area";
import {engines} from "../../wailsjs/go/models.ts";
import FileSelector from "@/components/FileSelector.vue";

import {ResizableHandle, ResizablePanel, ResizablePanelGroup,} from '@/components/ui/resizable'
import Settings from "@/components/Settings.vue";
import {RunPreprocessor} from "../../wailsjs/go/main/App";
import {useToast} from '@/components/ui/toast/use-toast'
import LanguageIcon from "@/components/icons/LanguageIcon.vue";
import RunArgs = engines.RunArgs;
import RunOptions = engines.RunOptions;

const {toast} = useToast()

const router = useRouter()

interface ConfigurationProps {
  engines: engines.EngineMeta[],
}

const selectedPath = ref("")

const preprocArgs = ref(<engines.EngineArg[]>[
      {
        name: "skip-empty",
        title: "Skip Empty Strings",
        description: "Do not obfuscate empty strings (eg. \"\")",
        type: "bool",
        defaultValue: true,
        value: "",
        availableOptions: <string[]>[],
        targetLanguage: ""
      },
      {
        name: "use-go-build-tag",
        title: "Use Go Build Tag",
        description: "Set a '//go:build <go tag>' constraint to process only targeted files",
        type: "bool",
        defaultValue: true,
        value: "",
        availableOptions: <string[]>[],
        targetLanguage: "golang"
      },
      {
        name: "go-build-tag",
        title: "Go Tag",
        description: "The value to set in the //go:build constraint to target specific files",
        type: "string",
        defaultValue: "obfuscate",
        value: "",
        availableOptions: <string[]>[],
        targetLanguage: "golang"
      },
      {
        name: "patch-marker",
        title: "Patch Marker",
        description: "The keyword used to mark strings as target for obfuscation (default: PATCH('mystring')",
        type: "string",
        defaultValue: "",
        value: "",
        availableOptions: <string[]>[],
        targetLanguage: "",
        notTargetLanguage: ""
      },
      {
        name: "target-extension",
        title: "Target Extension",
        description: "Filter the file that will be processed by extension",
        type: "[]string",
        defaultValue: [],
        value: [],
        availableOptions: <string[]>[],
        targetLanguage: "",
        notTargetLanguage: ""
      },
    ]
)

const preprocConfig = <RunOptions>{
  // targetDir: selectedPath.value,
}

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

const runPreproc = function () {
  if (selectedPath.value === "") {
    toast({
      title: "No path selected",
    })

    return
  }

  preprocArgs.value.forEach((arg) => {
    if (arg.name === "use-go-build-tag") {
      preprocConfig.useGoBuildTag = arg.value
    } else if (arg.name === "go-build-tag") {
      preprocConfig.goBuildTag = arg.value
    } else if (arg.name === "skip-empty") {
      preprocConfig.skipEmpty = arg.value
    } else if (arg.name === "patch-marker") {
      preprocConfig.patchMarker = arg.value
    } else if (arg.name === "target-extension") {
      preprocConfig.targetExtensions = arg.value
    }
  })

  preprocConfig.targetDir = selectedPath.value

  const args = <RunArgs>{
    engine: currentEngine.value.name.toLowerCase(),
    args: currentEngine.value.args,
    runOptions: preprocConfig
  }

  RunPreprocessor(args).then((res) => {
    toast({
      title: res,
    })

  }).catch((err) => {
    toast({
      title: 'Error running preprocessor' as string,
      description: err as string,
      variant: "destructive"
    })
  })

}
</script>

<template>
  <div class="flex-1 gap-4 overflow-auto p-4 px-6 lg:pt-10 lg:px-20">
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
        :args="currentEngine?.args"
        :current-engine="currentEngine"
        lg-cols="2"
        md-cols="2"
        sm-cols="1"
    />
  </div>

  <div class="flex flex-col sticky top-0 h-[calc(100vh_-_57px)] sm:w-80 lg:w-96 border-l lg:pt-10 p-4">
    <ResizablePanelGroup direction="vertical">
      <ResizablePanel :default-size="25">
        <!--          <div class="flex h-full items-center justify-center p-6">-->
        <!--            <span class="font-semibold">Two</span>-->
        <!--          </div>-->
        <ScrollArea class="flex-1 h-[10vh]">
          <h2 class="text-sm font-medium mb-4">Summary</h2>
          <div class="grid gap-4">
            <div
                v-for="(arg, _) in currentEngine.args as engines.EngineArg[]" :key="arg.name"
                class="grid grid-cols-2 items-center gap-4"
            >
              <Label for="arg-value" class=" font-semibold">
                {{ arg.title }}:
              </Label>
              <p id="arg-value"
                 :class="['text-sm break-all', arg.value === arg.defaultValue ? 'italic text-muted-foreground' : 'text-accent-foreground']">
                {{ arg.value === '' ? 'None' : arg.value }}
              </p>
            </div>
          </div>
        </ScrollArea>
      </ResizablePanel>
      <ResizableHandle with-handle/>

      <ResizablePanel :default-size="50">
        <ScrollArea class="py-4 h-full">
          <Settings
              :current-engine="currentEngine"
              :args="preprocArgs"
              legend="Settings"
          />
        </ScrollArea>
      </ResizablePanel>
      <ResizableHandle with-handle/>

      <ResizablePanel :default-size="25">
        <ScrollArea class="py-2 h-full">
          <div class="flex-1 pt-4">
            <h2 class="text-sm font-medium mb-4">Target</h2>
            <div class="grid gap-4">
            </div>
            <FileSelector
                @selected="(p) => {selectedPath = p}"
            />

          </div>
        </ScrollArea>

        <!--            <div class="flex h-full items-center justify-center p-6">-->
        <!--              <span class="font-semibold">Three</span>-->
        <!--            </div>-->
      </ResizablePanel>
    </ResizablePanelGroup>

    <!--    <Separator-->
    <!--        class="my-4"-->
    <!--    />-->

    <Separator
        class="my-4"
    />

    <div class="mt-auto">
      <Button
          class="w-full"
          size="lg"
          @click="runPreproc"
      >
        <Play/>
        Run
      </Button>
    </div>
  </div>
</template>

<style scoped>

</style>