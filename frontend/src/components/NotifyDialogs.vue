<script setup lang="ts">
import { computed, ref, shallowRef, useId, watch } from 'vue'
import type { Component } from 'vue'
import { CircleHelpIcon, InfoIcon, TriangleAlertIcon } from '@lucide/vue'
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogMedia,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import {
  Dialog as DialogRoot,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Field, FieldLabel } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { activeDialog } from '@/lib/notify'
import type { Dialog, NotifyIcon } from '@/lib/notify'

const icons: Record<NotifyIcon, Component> = {
  warning: TriangleAlertIcon,
  info: InfoIcon,
  question: CircleHelpIcon,
}

const confirmDialog = shallowRef<Dialog | null>(null)
const promptDialog = shallowRef<Dialog | null>(null)
const value = ref('')
const inputId = useId()

const confirmOpen = computed(() => activeDialog.value?.kind === 'confirm')
const promptOpen = computed(() => activeDialog.value?.kind === 'prompt')

watch(activeDialog, (d) => {
  if (d?.kind === 'confirm') confirmDialog.value = d
  if (d?.kind === 'prompt') {
    promptDialog.value = d
    value.value = d.inputValue
  }
})

function dismiss() {
  activeDialog.value?.close({ confirmed: false, value: '' })
}

function accept() {
  activeDialog.value?.close({ confirmed: true, value: value.value })
}

function onOpenChange(open: boolean) {
  if (!open) dismiss()
}

function focusInput(e: Event) {
  const input = document.getElementById(inputId)
  if (!(input instanceof HTMLInputElement)) return
  e.preventDefault()
  input.focus()
  input.select()
}
</script>

<template>
  <AlertDialog :open="confirmOpen" @update:open="onOpenChange">
    <AlertDialogContent v-if="confirmDialog">
      <AlertDialogHeader>
        <AlertDialogMedia
          v-if="confirmDialog.icon"
          :class="['size-12', confirmDialog.danger ? 'bg-destructive/10 text-destructive' : '']"
        >
          <component :is="icons[confirmDialog.icon]" class="size-6" />
        </AlertDialogMedia>
        <AlertDialogTitle>{{ confirmDialog.title }}</AlertDialogTitle>
        <AlertDialogDescription v-if="confirmDialog.text">{{ confirmDialog.text }}</AlertDialogDescription>
      </AlertDialogHeader>
      <AlertDialogFooter>
        <AlertDialogCancel>{{ confirmDialog.cancelText }}</AlertDialogCancel>
        <Button :variant="confirmDialog.danger ? 'destructive' : 'default'" @click="accept">
          {{ confirmDialog.confirmText }}
        </Button>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>

  <DialogRoot :open="promptOpen" @update:open="onOpenChange">
    <DialogContent v-if="promptDialog" class="sm:max-w-sm" @open-auto-focus="focusInput">
      <form class="grid gap-6" @submit.prevent="accept">
        <DialogHeader>
          <DialogTitle>{{ promptDialog.title }}</DialogTitle>
          <DialogDescription v-if="promptDialog.text">{{ promptDialog.text }}</DialogDescription>
        </DialogHeader>
        <Field>
          <FieldLabel v-if="promptDialog.inputLabel" :for="inputId">{{ promptDialog.inputLabel }}</FieldLabel>
          <Input :id="inputId" v-model="value" :placeholder="promptDialog.placeholder" autocomplete="off" />
        </Field>
        <DialogFooter>
          <Button type="button" variant="outline" @click="dismiss">{{ promptDialog.cancelText }}</Button>
          <Button type="submit">{{ promptDialog.confirmText }}</Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </DialogRoot>
</template>
