<template>
  <q-page class="q-pa-md q-pb-xl">
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h6 text-weight-bold text-primary">Admin - Match Results</div>
      <div class="row items-center q-gutter-sm">
        <q-btn
          label="Fetch Results"
          color="primary"
          icon="cloud_download"
          unelevated
          :loading="loading"
          @click="fetchExternalMatches"
        />
        <q-btn
          label="Sync Selected"
          color="positive"
          icon="sync"
          unelevated
          :loading="syncing"
          :disabled="selectedRows.length === 0"
          @click="syncSelected"
        />
      </div>
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
          { label: 'Pending', value: 'pending' },
          { label: 'Changed', value: 'changed' },
        ]"
      />
      <q-btn
        flat
        dense
        label="Select All Finished"
        icon="checklist"
        color="secondary"
        @click="selectAllFinished"
      />
      <q-space />
      <span v-if="lastSynced" class="text-caption text-grey-6">
        Last synced: {{ lastSynced }}
      </span>
    </div>

    <div v-if="loading" class="row justify-center q-my-xl">
      <q-spinner-dots color="primary" size="40px" />
    </div>

    <div v-else-if="matches.length === 0" class="text-center q-my-xl text-grey-6">
      No matches loaded. Click "Fetch Results" to load data from the external API.
    </div>

    <q-table
      v-else
      v-model:selected="selectedRows"
      :rows="filteredMatches"
      :columns="columns"
      row-key="match_number"
      selection="multiple"
      dense
      flat
      bordered
      :rows-per-page-options="[20, 50, 100]"
      :class="$q.dark.isActive ? 'bg-grey-10' : 'bg-white'"
    >
      <template #body-cell-home_team="props">
        <q-td :props="props">
          <span class="text-weight-medium">{{ props.row.external.home_team_name_en || props.row.local.home_team_name || props.row.external.home_team_label || 'TBD' }}</span>
        </q-td>
      </template>

      <template #body-cell-away_team="props">
        <q-td :props="props">
          <span class="text-weight-medium">{{ props.row.external.away_team_name_en || props.row.local.away_team_name || props.row.external.away_team_label || 'TBD' }}</span>
        </q-td>
      </template>

      <template #body-cell-score_home="props">
        <q-td :props="props">
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
      </template>

      <template #body-cell-score_away="props">
        <q-td :props="props">
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
      </template>

      <template #body-cell-status="props">
        <q-td :props="props">
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
      </template>

      <template #body-cell-ext_status="props">
        <q-td :props="props">
          <q-badge :color="getExtBadgeColor(props.row)">
            {{ getExtStatusText(props.row) }}
          </q-badge>
        </q-td>
      </template>

      <template #body-cell-diff="props">
        <q-td :props="props">
          <q-icon :name="getDiffIcon(props.row)" :color="getDiffColor(props.row)" size="sm" />
        </q-td>
      </template>

      <template #body="props">
        <q-tr
          :props="props"
          :class="getRowClass(props.row)"
        >
          <q-td><q-checkbox v-model="props.selected" dense /></q-td>
          <q-td key="match_number" :props="props">{{ props.row.match_number }}</q-td>
          <q-td key="group" :props="props">{{ props.row.group_code || props.row.phase }}</q-td>
          <q-td key="home_team" :props="props">
            <span class="text-weight-medium">{{ props.row.external.home_team_name_en || props.row.local.home_team_name || props.row.external.home_team_label || 'TBD' }}</span>
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
            <span class="text-weight-medium">{{ props.row.external.away_team_name_en || props.row.local.away_team_name || props.row.external.away_team_label || 'TBD' }}</span>
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
          <q-td key="ext_status" :props="props">
            <q-badge :color="getExtBadgeColor(props.row)">
              {{ getExtStatusText(props.row) }}
            </q-badge>
          </q-td>
          <q-td key="diff" :props="props">
            <q-icon :name="getDiffIcon(props.row)" :color="getDiffColor(props.row)" size="sm" />
          </q-td>
          <q-td key="actions" :props="props">
            <q-btn
              v-if="props.row.isDirty || (isExtFinished(props.row) && props.row.editStatus !== 'finished')"
              dense
              round
              flat
              icon="sync"
              color="positive"
              :loading="syncingRowNumber === props.row.match_number"
              @click="syncSingleRow(props.row)"
            >
              <q-tooltip>Sync this match</q-tooltip>
            </q-btn>
          </q-td>
        </q-tr>
      </template>
    </q-table>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { Notify } from 'quasar';
import { pb, PB_URL } from '@/boot/pocketbase';

interface ExternalMatchRow {
  match_number: number;
  phase: string;
  group_code: string;
  external: Record<string, unknown>;
  local: Record<string, unknown>;
  editScoreHome: number | null;
  editScoreAway: number | null;
  editStatus: string;
  isDirty: boolean;
}

const loading = ref(false);
const syncing = ref(false);
const syncingRowNumber = ref<number | null>(null);
const matches = ref<ExternalMatchRow[]>([]);
const selectedRows = ref<ExternalMatchRow[]>([]);
const selectedPhase = ref('all');
const filterMode = ref('all');
const lastSynced = ref('');

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
  { name: 'ext_status', label: 'Ext', field: 'ext_status', align: 'center' as const },
  { name: 'diff', label: 'Diff', field: 'diff', align: 'center' as const },
  { name: 'actions', label: '', field: 'actions', align: 'center' as const },
];

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
    result = result.filter((m) => isExtFinished(m));
  } else if (filterMode.value === 'pending') {
    result = result.filter((m) => !isExtFinished(m));
  } else if (filterMode.value === 'changed') {
    result = result.filter((m) => isChanged(m));
  }

  return result;
});

function isExtFinished(row: ExternalMatchRow): boolean {
  const finished = row.external.finished;
  if (typeof finished === 'string') return finished.toUpperCase() === 'TRUE';
  return finished === true;
}

function isChanged(row: ExternalMatchRow): boolean {
  return row.isDirty;
}

function getExtStatusText(row: ExternalMatchRow): string {
  if (isExtFinished(row)) return 'Done';
  const elapsed = row.external.time_elapsed;
  if (typeof elapsed === 'string' && elapsed !== 'notstarted') return elapsed;
  return 'Pending';
}

function getExtBadgeColor(row: ExternalMatchRow): string {
  if (isExtFinished(row)) return 'positive';
  const elapsed = row.external.time_elapsed;
  if (typeof elapsed === 'string' && elapsed !== 'notstarted') return 'warning';
  return 'grey';
}

function getDiffIcon(row: ExternalMatchRow): string {
  if (row.isDirty) return 'edit';
  if (isExtFinished(row) && row.editStatus === 'finished') return 'check_circle';
  if (!isExtFinished(row)) return 'hourglass_empty';
  return 'warning';
}

function getDiffColor(row: ExternalMatchRow): string {
  if (row.isDirty) return 'warning';
  if (isExtFinished(row) && row.editStatus === 'finished') return 'positive';
  if (!isExtFinished(row)) return 'grey';
  return 'warning';
}

function getRowClass(row: ExternalMatchRow): string {
  if (row.isDirty) return 'bg-orange-1';
  if (isExtFinished(row) && row.editStatus === 'finished') return $q_dark() ? 'bg-green-10' : 'bg-green-1';
  return '';
}

function $q_dark(): boolean {
  return document.body.classList.contains('body--dark');
}

function markDirty(row: ExternalMatchRow) {
  const localStatus = row.local.status as string;

  const isSameAsOriginal =
    row.editScoreHome === (row.local.score_home as number | null) &&
    row.editScoreAway === (row.local.score_away as number | null) &&
    row.editStatus === localStatus;

  row.isDirty = !isSameAsOriginal;
}

function parseExtScore(val: unknown): number | null {
  if (val === null || val === undefined) return null;
  const n = typeof val === 'string' ? parseInt(val, 10) : Number(val);
  return isNaN(n) ? null : n;
}

function selectAllFinished() {
  const finishedRows = filteredMatches.value.filter((m) => isExtFinished(m));
  selectedRows.value = finishedRows;
}

async function fetchExternalMatches() {
  loading.value = true;
  try {
    const response = await fetch(`${PB_URL}/api/wc/external/matches`, {
      headers: {
        Authorization: `Bearer ${pb.authStore.token}`,
      },
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.message || `HTTP ${response.status}`);
    }

    const data = await response.json();
    matches.value = data.map((item: ExternalMatchRow) => {
      const extFinished = isExtFinishedRaw(item.external);
      const extHome = extFinished ? parseExtScore(item.external.home_score) : null;
      const extAway = extFinished ? parseExtScore(item.external.away_score) : null;

      return {
        ...item,
        editScoreHome: item.local.score_home != null ? item.local.score_home as number : extHome,
        editScoreAway: item.local.score_away != null ? item.local.score_away as number : extAway,
        editStatus: (item.local.status as string) || (extFinished ? 'finished' : 'upcoming'),
        isDirty: false,
      };
    });

    Notify.create({ type: 'positive', message: `Loaded ${data.length} matches from external API` });
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Unknown error';
    Notify.create({ type: 'negative', message: `Failed to fetch: ${message}` });
  } finally {
    loading.value = false;
  }
}

function isExtFinishedRaw(external: Record<string, unknown>): boolean {
  const finished = external.finished;
  if (typeof finished === 'string') return finished.toUpperCase() === 'TRUE';
  return finished === true;
}

async function syncSelected() {
  if (selectedRows.value.length === 0) return;

  syncing.value = true;
  try {
    const payload = selectedRows.value.map((row) => ({
      match_number: row.match_number,
      score_home: row.editScoreHome ?? 0,
      score_away: row.editScoreAway ?? 0,
      status: row.editStatus,
    }));

    const response = await fetch(`${PB_URL}/api/wc/external/sync`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${pb.authStore.token}`,
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.message || `HTTP ${response.status}`);
    }

    const result = await response.json();
    Notify.create({ type: 'positive', message: result.message || 'Synced successfully' });

    lastSynced.value = new Date().toLocaleTimeString();

    for (const row of selectedRows.value) {
      row.local.score_home = row.editScoreHome;
      row.local.score_away = row.editScoreAway;
      row.local.status = row.editStatus;
      row.isDirty = false;
    }
    selectedRows.value = [];
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Unknown error';
    Notify.create({ type: 'negative', message: `Sync failed: ${message}` });
  } finally {
    syncing.value = false;
  }
}

async function syncSingleRow(row: ExternalMatchRow) {
  syncingRowNumber.value = row.match_number;
  try {
    const payload = [
      {
        match_number: row.match_number,
        score_home: row.editScoreHome ?? 0,
        score_away: row.editScoreAway ?? 0,
        status: row.editStatus,
      },
    ];

    const response = await fetch(`${PB_URL}/api/wc/external/sync`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${pb.authStore.token}`,
      },
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      const err = await response.json().catch(() => ({}));
      throw new Error(err.message || `HTTP ${response.status}`);
    }

    row.local.score_home = row.editScoreHome;
    row.local.score_away = row.editScoreAway;
    row.local.status = row.editStatus;
    row.isDirty = false;

    Notify.create({ type: 'positive', message: `Match #${row.match_number} synced` });
  } catch (err: unknown) {
    const message = err instanceof Error ? err.message : 'Unknown error';
    Notify.create({ type: 'negative', message: `Sync failed: ${message}` });
  } finally {
    syncingRowNumber.value = null;
  }
}
</script>

<style scoped>
.bg-orange-1 {
  background-color: rgba(255, 152, 0, 0.08);
}
</style>
