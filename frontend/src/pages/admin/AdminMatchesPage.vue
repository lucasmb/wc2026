<template>
  <q-page class="q-pa-md q-pb-xl">
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h6 text-weight-bold text-primary">Admin - Edit Matches</div>
      <q-btn
        label="Save Changes"
        color="positive"
        icon="save"
        unelevated
        :loading="saving"
        :disabled="!hasChanges"
        @click="saveChanges"
      />
    </div>

    <div class="row items-center q-gutter-sm q-mb-md">
      <q-select
        v-model="selectedPhase"
        :options="phaseOptions"
        option-value="value"
        option-label="label"
        emit-value
        map-options
        outlined
        dense
        label="Phase"
        style="width: 180px"
      />
      <q-btn-toggle
        v-model="filterMode"
        toggle-color="primary"
        no-caps
        dense
        :options="[
          { label: 'All', value: 'all' },
          { label: 'Finished', value: 'finished' },
          { label: 'Upcoming', value: 'upcoming' },
          { label: 'Changed', value: 'changed' },
        ]"
      />
    </div>

    <div v-if="loading" class="row justify-center q-my-xl">
      <q-spinner-dots color="primary" size="40px" />
    </div>

    <q-table
      v-else
      :rows="filteredMatches"
      :columns="columns"
      row-key="id"
      dense
      flat
      bordered
      :rows-per-page-options="[20, 50, 100]"
      :class="$q.dark.isActive ? 'bg-grey-10' : 'bg-white'"
    >
      <template #body="props">
        <q-tr
          :props="props"
          :class="props.row.isDirty ? ($q.dark.isActive ? 'row-dirty-dark' : 'row-dirty-light') : ''"
        >
          <q-td key="match_number" :props="props">{{ props.row.match_number }}</q-td>
          <q-td key="group" :props="props">{{ props.row.group_code || props.row.phase }}</q-td>
          <q-td key="home_team" :props="props">
            <span class="text-weight-medium">{{ props.row.expand?.home_team?.name || 'TBD' }}</span>
          </q-td>
          <q-td key="score_home" :props="props">
            <q-input
              v-model.number="props.row.editScoreHome"
              type="number"
              dense
              outlined
              style="width: 60px"
              input-class="text-center"
              @update:model-value="markDirty(props.row)"
            />
          </q-td>
          <q-td key="score_away" :props="props">
            <q-input
              v-model.number="props.row.editScoreAway"
              type="number"
              dense
              outlined
              style="width: 60px"
              input-class="text-center"
              @update:model-value="markDirty(props.row)"
            />
          </q-td>
          <q-td key="away_team" :props="props">
            <span class="text-weight-medium">{{ props.row.expand?.away_team?.name || 'TBD' }}</span>
          </q-td>
          <q-td key="status" :props="props">
            <q-select
              v-model="props.row.editStatus"
              :options="statusOptions"
              dense
              outlined
              emit-value
              map-options
              style="width: 110px"
              @update:model-value="markDirty(props.row)"
            />
          </q-td>
          <q-td key="winner" :props="props">
            <span v-if="props.row.editStatus === 'finished' && props.row.editScoreHome !== null && props.row.editScoreAway !== null">
              {{ getWinnerName(props.row) }}
            </span>
            <span v-else class="text-grey-5">-</span>
          </q-td>
          <q-td key="actions" :props="props">
            <q-btn
              v-if="props.row.isDirty"
              dense
              round
              flat
              icon="save"
              color="positive"
              :loading="savingRowId === props.row.id"
              @click="saveSingleRow(props.row)"
            >
              <q-tooltip>Save this match</q-tooltip>
            </q-btn>
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { Notify } from 'quasar';
import { pb } from '@/boot/pocketbase';
import { useTournamentStore } from '@/stores/tournament';
import type { Match } from '@/types';

interface EditableMatch extends Match {
  editScoreHome: number | null;
  editScoreAway: number | null;
  editStatus: string;
  isDirty: boolean;
}

const tournamentStore = useTournamentStore();
const loading = ref(true);
const saving = ref(false);
const savingRowId = ref<string | null>(null);
const matches = ref<EditableMatch[]>([]);
const selectedPhase = ref('all');
const filterMode = ref('all');

const phaseOptions = [
  { label: 'All', value: 'all' },
  { label: 'Group A', value: 'A' },
  { label: 'Group B', value: 'B' },
  { label: 'Group C', value: 'C' },
  { label: 'Group D', value: 'D' },
  { label: 'Group E', value: 'E' },
  { label: 'Group F', value: 'F' },
  { label: 'Group G', value: 'G' },
  { label: 'Group H', value: 'H' },
  { label: 'Group I', value: 'I' },
  { label: 'Group J', value: 'J' },
  { label: 'Group K', value: 'K' },
  { label: 'Group L', value: 'L' },
  { label: 'Round of 32', value: 'r32' },
  { label: 'Round of 16', value: 'r16' },
  { label: 'Quarter-Finals', value: 'qf' },
  { label: 'Semi-Finals', value: 'sf' },
  { label: 'Third Place', value: 'third' },
  { label: 'Final', value: 'final' },
];

const statusOptions = [
  { label: 'Upcoming', value: 'upcoming' },
  { label: 'Live', value: 'live' },
  { label: 'Finished', value: 'finished' },
];

const columns = [
  { name: 'match_number', label: '#', field: 'match_number', align: 'center' as const, sortable: true },
  { name: 'group', label: 'Group', field: 'group_code', align: 'center' as const, sortable: true },
  { name: 'home_team', label: 'Home', field: 'home_team', align: 'left' as const },
  { name: 'score_home', label: 'H', field: 'score_home', align: 'center' as const },
  { name: 'score_away', label: 'A', field: 'score_away', align: 'center' as const },
  { name: 'away_team', label: 'Away', field: 'away_team', align: 'left' as const },
  { name: 'status', label: 'Status', field: 'status', align: 'center' as const },
  { name: 'winner', label: 'Winner', field: 'winner', align: 'left' as const },
  { name: 'actions', label: '', field: 'actions', align: 'center' as const },
];

const hasChanges = computed(() => matches.value.some((m) => m.isDirty));

const filteredMatches = computed(() => {
  let result = matches.value;

  if (selectedPhase.value !== 'all') {
    const phaseLetters = ['r32', 'r16', 'qf', 'sf', 'third', 'final'];
    if (phaseLetters.includes(selectedPhase.value)) {
      result = result.filter((m) => m.phase === selectedPhase.value);
    } else {
      result = result.filter((m) => m.phase === 'group' && m.group_code === selectedPhase.value);
    }
  }

  if (filterMode.value === 'finished') {
    result = result.filter((m) => m.editStatus === 'finished');
  } else if (filterMode.value === 'upcoming') {
    result = result.filter((m) => m.editStatus === 'upcoming');
  } else if (filterMode.value === 'changed') {
    result = result.filter((m) => m.isDirty);
  }

  return result;
});

function markDirty(row: EditableMatch) {
  const isSameAsOriginal =
    row.editScoreHome === (row.score_home ?? null) &&
    row.editScoreAway === (row.score_away ?? null) &&
    row.editStatus === row.status;

  row.isDirty = !isSameAsOriginal;
}

function getWinnerName(row: EditableMatch): string {
  if (row.editScoreHome === null || row.editScoreAway === null) return '-';
  if (row.editScoreHome > row.editScoreAway) {
    return row.expand?.home_team?.name || 'Home';
  } else if (row.editScoreAway > row.editScoreHome) {
    return row.expand?.away_team?.name || 'Away';
  }
  return 'Draw';
}

async function loadMatches() {
  loading.value = true;
  try {
    await tournamentStore.fetchMatches();
    matches.value = tournamentStore.matches.map((m) => ({
      ...m,
      editScoreHome: m.score_home ?? null,
      editScoreAway: m.score_away ?? null,
      editStatus: m.status,
      isDirty: false,
    }));
  } finally {
    loading.value = false;
  }
}

async function saveChanges() {
  const dirtyMatches = matches.value.filter((m) => m.isDirty);
  if (dirtyMatches.length === 0) return;

  saving.value = true;
  try {
    for (const match of dirtyMatches) {
      const updateData: Record<string, unknown> = {
        status: match.editStatus,
      };

      if (match.editStatus === 'finished' || match.editStatus === 'live') {
        updateData.score_home = match.editScoreHome ?? 0;
        updateData.score_away = match.editScoreAway ?? 0;

        if (match.editStatus === 'finished') {
          const homeScore = match.editScoreHome ?? 0;
          const awayScore = match.editScoreAway ?? 0;
          if (homeScore > awayScore) {
            updateData.winner = match.home_team;
          } else if (awayScore > homeScore) {
            updateData.winner = match.away_team;
          } else {
            updateData.winner = null;
          }
        }
      } else {
        updateData.score_home = null;
        updateData.score_away = null;
        updateData.winner = null;
      }

      await pb.collection('matches_id').update(match.id, updateData);
      if (match.editScoreHome !== null) {
        match.score_home = match.editScoreHome;
      } else {
        delete match.score_home;
      }
      if (match.editScoreAway !== null) {
        match.score_away = match.editScoreAway;
      } else {
        delete match.score_away;
      }
      match.status = match.editStatus as Match['status'];
      match.isDirty = false;
    }

    Notify.create({ type: 'positive', message: `Saved ${dirtyMatches.length} match(es)` });
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Unknown error';
    Notify.create({ type: 'negative', message: `Save failed: ${message}` });
  } finally {
    saving.value = false;
  }
}

async function saveSingleRow(match: EditableMatch) {
  savingRowId.value = match.id;
  try {
    const updateData: Record<string, unknown> = {
      status: match.editStatus,
    };

    if (match.editStatus === 'finished' || match.editStatus === 'live') {
      updateData.score_home = match.editScoreHome ?? 0;
      updateData.score_away = match.editScoreAway ?? 0;

      if (match.editStatus === 'finished') {
        const homeScore = match.editScoreHome ?? 0;
        const awayScore = match.editScoreAway ?? 0;
        if (homeScore > awayScore) {
          updateData.winner = match.home_team;
        } else if (awayScore > homeScore) {
          updateData.winner = match.away_team;
        } else {
          updateData.winner = null;
        }
      }
    } else {
      updateData.score_home = null;
      updateData.score_away = null;
      updateData.winner = null;
    }

    await pb.collection('matches_id').update(match.id, updateData);
    if (match.editScoreHome !== null) {
      match.score_home = match.editScoreHome;
    } else {
      delete match.score_home;
    }
    if (match.editScoreAway !== null) {
      match.score_away = match.editScoreAway;
    } else {
      delete match.score_away;
    }
    match.status = match.editStatus as Match['status'];
    match.isDirty = false;

    Notify.create({ type: 'positive', message: `Match #${match.match_number} saved` });
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Unknown error';
    Notify.create({ type: 'negative', message: `Save failed: ${message}` });
  } finally {
    savingRowId.value = null;
  }
}

onMounted(() => {
  void loadMatches();
});
</script>

<style scoped>
.row-dirty-light {
  background-color: rgba(255, 152, 0, 0.12);
}
.row-dirty-dark {
  background-color: rgba(255, 152, 0, 0.25);
}
</style>
