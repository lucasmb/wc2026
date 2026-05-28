<template>
  <q-card flat bordered class="rounded-borders q-pa-sm bg-white">
    <div class="row items-center justify-around q-col-gutter-sm">
      <!-- Home Side -->
      <div class="col-4 text-center">
        <q-img
          v-if="match.expand?.home_team?.flag_url"
          :src="match.expand.home_team.flag_url"
          style="
            width: 48px;
            height: 32px;
            border-radius: 4px;
            border: 1px solid rgba(0, 0, 0, 0.1);
          "
        />
        <div class="text-subtitle2 text-weight-bold text-grey-9 q-mt-xs text-ellipsis">
          {{ match.expand?.home_team?.name }}
        </div>
      </div>

      <!-- Score Inputs (or locking placeholders) -->
      <div class="col-4 text-center">
        <div v-if="isLocked" class="row items-center justify-center q-gutter-x-sm">
          <span class="text-h5 text-weight-bold text-grey-7">{{ form.predicted_home }}</span>
          <span class="text-grey-5">:</span>
          <span class="text-h5 text-weight-bold text-grey-7">{{ form.predicted_away }}</span>
        </div>
        <div v-else class="row items-center justify-center q-gutter-x-sm">
          <q-input
            v-model.number="form.predicted_home"
            type="number"
            min="0"
            outlined
            dense
            input-class="text-center text-weight-bold text-subtitle1 no-spin"
            style="width: 60px"
            hide-bottom-space
          />
          <span class="text-grey-7 text-weight-bold">:</span>
          <q-input
            v-model.number="form.predicted_away"
            type="number"
            min="0"
            outlined
            dense
            input-class="text-center text-weight-bold text-subtitle1 no-spin"
            style="width: 60px"
            hide-bottom-space
          />
        </div>
      </div>

      <!-- Away Side -->
      <div class="col-4 text-center">
        <q-img
          v-if="match.expand?.away_team?.flag_url"
          :src="match.expand.away_team.flag_url"
          style="
            width: 48px;
            height: 32px;
            border-radius: 4px;
            border: 1px solid rgba(0, 0, 0, 0.1);
          "
        />
        <div class="text-subtitle2 text-weight-bold text-grey-9 q-mt-xs text-ellipsis">
          {{ match.expand?.away_team?.name }}
        </div>
      </div>
    </div>

    <q-separator class="q-my-md" />

    <div class="row items-center justify-between">
      <div>
        <q-chip v-if="isLocked" dense color="grey-3" text-color="grey-8" icon="lock">
          Predictions Locked
        </q-chip>
        <q-chip v-else dense color="green-1" text-color="green-8" icon="schedule"> Open </q-chip>

        <!-- Display Points Awarded if match has finished -->
        <q-chip
          v-if="match.status === 'finished' && predictionId"
          dense
          color="blue-1"
          text-color="blue-8"
          icon="emoji_events"
          class="text-weight-bold"
        >
          +{{ pointsAwarded }} Pts
        </q-chip>
      </div>

      <q-btn
        v-if="!isLocked"
        label="Save"
        color="primary"
        dense
        unelevated
        class="q-px-md"
        :loading="saving"
        @click="savePrediction"
      />
    </div>
  </q-card>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { Notify } from 'quasar';
import { pb } from '@/boot/pocketbase';
import { useAuthStore } from '@/stores/auth';
import type { Match, Prediction } from '@/types';

const props = defineProps<{
  match: Match;
  groupId: string;
  prediction: Prediction | undefined;
}>();

const emit = defineEmits<{
  (e: 'saved', prediction: Prediction): void;
}>();

const authStore = useAuthStore();
const saving = ref(false);
const predictionId = ref<string | null>(null);
const pointsAwarded = ref(0);

const form = ref({
  predicted_home: 0,
  predicted_away: 0,
});

const isLocked = computed(() => {
  if (props.match.status === 'finished') return true;
  if (!props.match.kickoff) return false;
  const isoString = props.match.kickoff.replace(' ', 'T');
  return new Date(isoString).getTime() <= Date.now();
});

// Reactively bind incoming batched properties directly from prop watchers
watch(
  () => props.prediction,
  (newPred) => {
    if (newPred) {
      predictionId.value = newPred.id;
      form.value.predicted_home = newPred.predicted_home;
      form.value.predicted_away = newPred.predicted_away;
      pointsAwarded.value = newPred.points_awarded || 0;
    } else {
      predictionId.value = null;
      form.value.predicted_home = 0;
      form.value.predicted_away = 0;
      pointsAwarded.value = 0;
    }
  },
  { immediate: true },
);

async function savePrediction() {
  if (isLocked.value || !authStore.user?.id || !props.groupId) return;
  saving.value = true;
  try {
    const data = {
      user: authStore.user.id,
      match: props.match.id,
      prediction_group: props.groupId,
      predicted_home: Math.max(0, form.value.predicted_home),
      predicted_away: Math.max(0, form.value.predicted_away),
    };

    let resultRecord: Prediction;

    if (predictionId.value) {
      const updated = await pb.collection('predictions_id').update(predictionId.value, data);
      resultRecord = updated as unknown as Prediction;
    } else {
      const created = await pb.collection('predictions_id').create(data);
      resultRecord = created as unknown as Prediction;
    }

    // Emit the saved result upward to update the parent’s central batch map cache
    emit('saved', resultRecord);

    Notify.create({
      type: 'positive',
      message: 'Prediction saved successfully!',
    });
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Check your internet connection';
    Notify.create({
      type: 'negative',
      message: `Failed saving prediction: ${message}`,
    });
  } finally {
    saving.value = false;
  }
}
</script>

<style scoped>
.text-ellipsis {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

:deep(input::-webkit-outer-spin-button),
:deep(input::-webkit-inner-spin-button) {
  -webkit-appearance: none;
  margin: 0;
}

:deep(input[type='number']) {
  -moz-appearance: textfield;
}
</style>
