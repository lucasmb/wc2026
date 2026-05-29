<template>
  <!-- Dynamic page background class -->
  <q-page class="q-pa-md" :class="$q.dark.isActive ? 'bg-dark' : 'bg-grey-1'">
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h6 text-weight-bold text-primary">Llaves y Partidos del Torneo</div>

      <!-- Selector de Grupo de Predicción (Se mantiene alineado con watchers reactivos) -->
      <q-select
        v-if="myGroups.length > 0"
        v-model="activeGroupId"
        :options="myGroups"
        :dark="$q.dark.isActive"
        option-value="id"
        option-label="name"
        emit-value
        map-options
        label="Seleccionar Grupo"
        outlined
        dense
        style="width: 190px"
      />
    </div>

    <!-- Mode Tabs -->
    <q-tabs
      v-model="activeTab"
      dense
      :dark="$q.dark.isActive"
      :class="$q.dark.isActive ? 'bg-grey-10' : 'bg-white'"
      class="text-grey rounded-borders shadow-1"
      active-color="primary"
      indicator-color="primary"
      align="justify"
    >
      <q-tab name="fixtures" icon="list" label="Mis Predicciones" />
      <q-tab name="bracket" icon="account_tree" label="Llaves de Eliminación" />
    </q-tabs>

    <q-tab-panels v-model="activeTab" animated class="bg-transparent q-mt-md">
      <!-- PESTAÑA 1: Listado de Partidos y Predicciones -->
      <q-tab-panel name="fixtures" class="q-pa-none">
        <div v-if="loading" class="row justify-center q-my-xl">
          <q-spinner-dots color="primary" size="40px" />
        </div>

        <q-card
          v-else
          flat
          bordered
          :dark="$q.dark.isActive"
          :class="$q.dark.isActive ? 'bg-grey-10 border-grey-8' : 'bg-white'"
          class="rounded-borders shadow-1"
        >
          <q-card-section class="q-pb-none row items-center justify-between q-col-gutter-sm">
            <div class="text-subtitle1 text-weight-bold text-primary">
              Todos los Partidos e Historial
            </div>

            <!-- Controles de Filtros -->
            <div class="row q-gutter-x-sm">
              <!-- Filtro de Fase -->
              <q-select
                v-model="selectedPhaseFilter"
                :options="phaseFilterOptions"
                :dark="$q.dark.isActive"
                label="Filtrar por Fase"
                outlined
                dense
                emit-value
                map-options
                style="width: 170px"
              />

              <!-- Búsqueda por Texto -->
              <q-input
                v-model="filterText"
                :dark="$q.dark.isActive"
                dense
                placeholder="Buscar países..."
                outlined
                style="width: 180px"
              >
                <template v-slot:append><q-icon name="search" /></template>
              </q-input>
            </div>
          </q-card-section>

          <q-card-section>
            <q-table
              flat
              bordered
              :dark="$q.dark.isActive"
              :rows="filteredFixtures"
              :columns="columns"
              row-key="id"
              v-model:pagination="pagination"
              class="no-shadow"
            >
              <!-- Slot de Columna de Fase Traducida -->
              <template v-slot:body-cell-phase="props">
                <q-td :props="props" class="text-center">
                  <q-badge color="blue-1" text-color="blue-8" class="text-weight-bold">
                    {{ translatePhase(props.row.phase) }}
                  </q-badge>
                </q-td>
              </template>

              <!-- Mapped Home Team Flag -->
              <template v-slot:body-cell-home_team="props">
                <q-td :props="props" class="row items-center q-gutter-x-sm">
                  <q-img
                    v-if="props.row.expand?.home_team?.flag_url"
                    :src="props.row.expand.home_team.flag_url"
                    style="width: 24px; height: 16px; border-radius: 2px"
                  />
                  <span class="text-weight-medium">{{ props.row.expand?.home_team?.name }}</span>
                </q-td>
              </template>

              <!-- Mapped Away Team Flag -->
              <template v-slot:body-cell-away_team="props">
                <q-td :props="props" class="row items-center q-gutter-x-sm">
                  <q-img
                    v-if="props.row.expand?.away_team?.flag_url"
                    :src="props.row.expand.away_team.flag_url"
                    style="width: 24px; height: 16px; border-radius: 2px"
                  />
                  <span class="text-weight-medium">{{ props.row.expand?.away_team?.name }}</span>
                </q-td>
              </template>

              <!-- Predicción del usuario -->
              <template v-slot:body-cell-prediction="props">
                <q-td :props="props" class="text-center">
                  <span v-if="predictionsMap[props.row.id]" class="text-weight-bold text-blue-5">
                    {{ predictionsMap[props.row.id]?.predicted_home }} -
                    {{ predictionsMap[props.row.id]?.predicted_away }}
                  </span>
                  <span v-else class="text-grey-5 text-weight-medium">Sin Predicción</span>
                </q-td>
              </template>

              <!-- Resultado Real del Partido -->
              <template v-slot:body-cell-actual_score="props">
                <q-td :props="props" class="text-center text-weight-bold">
                  <span v-if="props.row.status === 'finished'">
                    {{ props.row.score_home }} - {{ props.row.score_away }}
                  </span>
                  <span v-else class="text-grey-5 text-weight-regular">-</span>
                </q-td>
              </template>

              <!-- Puntos Ganados -->
              <template v-slot:body-cell-points="props">
                <q-td :props="props" class="text-center">
                  <q-badge
                    v-if="props.row.status === 'finished' && predictionsMap[props.row.id]"
                    :color="
                      (predictionsMap[props.row.id]?.points_awarded ?? 0) > 0
                        ? 'green-1'
                        : $q.dark.isActive
                          ? 'grey-9'
                          : 'grey-2'
                    "
                    :text-color="
                      (predictionsMap[props.row.id]?.points_awarded ?? 0) > 0
                        ? 'green-8'
                        : $q.dark.isActive
                          ? 'grey-4'
                          : 'grey-8'
                    "
                    label="Puntos"
                    class="text-weight-bold q-px-sm"
                  >
                    +{{ predictionsMap[props.row.id]?.points_awarded ?? 0 }}
                  </q-badge>
                  <span v-else>-</span>
                </q-td>
              </template>
            </q-table>
          </q-card-section>
        </q-card>
      </q-tab-panel>

      <!-- PESTAÑA 2: Llave de eliminación interactiva -->
      <q-tab-panel name="bracket" class="q-pa-none">
        <q-scroll-area style="height: 600px; width: 100%">
          <div class="row no-wrap q-gutter-x-lg q-pa-sm" style="min-width: 1400px">
            <!-- COLUMNA 1: 16avos de Final -->
            <div class="bracket-col">
              <div
                class="text-subtitle2 text-weight-bold text-center text-primary text-uppercase q-mb-md"
              >
                16avos de Final
              </div>
              <div class="q-gutter-y-md">
                <div
                  v-for="i in 8"
                  :key="'r32-' + i"
                  class="bracket-match border rounded-borders q-pa-sm"
                  :class="
                    $q.dark.isActive
                      ? 'bg-grey-9 border-grey-8 text-white'
                      : 'bg-grey-1 text-grey-8'
                  "
                >
                  <div class="row justify-between text-caption text-weight-bold">
                    <span>Clasificado {{ i * 2 - 1 }}</span>
                    <span>-</span>
                  </div>
                  <q-separator class="q-my-xs" :dark="$q.dark.isActive" />
                  <div class="row justify-between text-caption text-weight-bold">
                    <span>Clasificado {{ i * 2 }}</span>
                    <span>-</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- COLUMNA 2: Octavos de Final -->
            <div class="bracket-col justify-around column">
              <div
                class="text-subtitle2 text-weight-bold text-center text-primary text-uppercase q-mb-md"
              >
                Octavos de Final
              </div>
              <div class="q-gutter-y-xl">
                <div
                  v-for="i in 4"
                  :key="'r16-' + i"
                  class="bracket-match border rounded-borders q-pa-sm"
                  :class="
                    $q.dark.isActive
                      ? 'bg-grey-9 border-grey-8 text-white'
                      : 'bg-grey-1 text-grey-8'
                  "
                >
                  <div class="row justify-between text-caption text-weight-bold">
                    <span>Ganador Partido {{ i * 2 - 1 }}</span>
                    <span>-</span>
                  </div>
                  <q-separator class="q-my-xs" :dark="$q.dark.isActive" />
                  <div class="row justify-between text-caption text-weight-bold">
                    <span>Ganador Partido {{ i * 2 }}</span>
                    <span>-</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- COLUMNA 3: Cuartos de Final -->
            <div class="bracket-col justify-around column">
              <div
                class="text-subtitle2 text-weight-bold text-center text-primary text-uppercase q-mb-md"
              >
                Cuartos de Final
              </div>
              <div class="q-gutter-y-lg" style="margin-top: 50px">
                <div
                  v-for="i in 2"
                  :key="'qf-' + i"
                  class="bracket-match border rounded-borders q-pa-sm"
                  :class="
                    $q.dark.isActive
                      ? 'bg-grey-9 border-grey-8 text-white'
                      : 'bg-grey-1 text-grey-8'
                  "
                >
                  <div class="row justify-between text-caption text-weight-bold">
                    <span>Ganador Cuartos {{ i * 2 - 1 }}</span>
                    <span>-</span>
                  </div>
                  <q-separator class="q-my-xs" :dark="$q.dark.isActive" />
                  <div class="row justify-between text-caption text-weight-bold">
                    <span>Ganador Cuartos {{ i * 2 }}</span>
                    <span>-</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- COLUMNA 4: Semifinales -->
            <div class="bracket-col justify-center column">
              <div
                class="text-subtitle2 text-weight-bold text-center text-primary text-uppercase q-mb-md"
              >
                Semifinales
              </div>
              <div
                class="bracket-match border rounded-borders q-pa-sm"
                :class="
                  $q.dark.isActive ? 'bg-grey-9 border-grey-8 text-white' : 'bg-grey-1 text-grey-8'
                "
              >
                <div class="row justify-between text-caption text-weight-bold">
                  <span>Ganador Semifinal 1</span>
                  <span>-</span>
                </div>
                <q-separator class="q-my-xs" :dark="$q.dark.isActive" />
                <div class="row justify-between text-caption text-weight-bold">
                  <span>Ganador Semifinal 2</span>
                  <span>-</span>
                </div>
              </div>
            </div>

            <!-- COLUMNA 5: Final Mundial -->
            <div class="bracket-col justify-center column">
              <div
                class="text-subtitle2 text-weight-bold text-center text-primary text-uppercase q-mb-md"
              >
                Gran Final
              </div>
              <div
                class="bracket-match border rounded-borders q-pa-md bg-primary text-white shadow-3"
              >
                <div class="text-center text-caption text-weight-bolder text-uppercase q-mb-xs">
                  Final 19 de Julio
                </div>
                <div class="row justify-between text-caption text-weight-bold">
                  <span>Ganador Semifinal 1</span>
                  <span>-</span>
                </div>
                <q-separator dark class="q-my-xs" />
                <div class="row justify-between text-caption text-weight-bold">
                  <span>Ganador Semifinal 2</span>
                  <span>-</span>
                </div>
              </div>
            </div>
          </div>
        </q-scroll-area>
      </q-tab-panel>
    </q-tab-panels>
  </q-page>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';
import { pb } from '@/boot/pocketbase';
import { useAuthStore } from '@/stores/auth';
import { useTournamentStore } from '@/stores/tournament';
import type { PredictionGroup, Prediction } from '@/types';
import type { QTableColumn } from 'quasar';

const authStore = useAuthStore();
const tournamentStore = useTournamentStore();

const activeTab = ref<'fixtures' | 'bracket'>('fixtures');
const loading = ref(true);
const filterText = ref('');
const selectedPhaseFilter = ref('all');

const activeGroupId = ref('');
const myGroups = ref<PredictionGroup[]>([]);
const predictionsMap = ref<Record<string, Prediction>>({});

// Default Table sorting configuration (match_number ASC)
const pagination = ref({
  sortBy: 'match_number',
  descending: false,
  rowsPerPage: 15,
});

const phaseFilterOptions = [
  { label: 'Todos los partidos', value: 'all' },
  { label: 'Fase de Grupos', value: 'group' },
  { label: '16avos de Final', value: 'r32' },
  { label: 'Octavos de Final', value: 'r16' },
  { label: 'Cuartos de Final', value: 'qf' },
  { label: 'Semifinales', value: 'sf' },
  { label: 'Final', value: 'final' },
];

const columns: QTableColumn[] = [
  {
    name: 'match_number',
    align: 'center',
    label: 'Partido #',
    field: 'match_number',
    sortable: true,
  },
  { name: 'phase', align: 'center', label: 'Fase', field: 'phase', sortable: true },
  { name: 'group_code', align: 'center', label: 'Grupo', field: 'group_code', sortable: true },
  { name: 'home_team', align: 'left', label: 'Local', field: 'home_team' },
  { name: 'prediction', align: 'center', label: 'Mi Pred.', field: 'id' },
  { name: 'actual_score', align: 'center', label: 'Resultado', field: 'id' },
  { name: 'away_team', align: 'left', label: 'Visitante', field: 'away_team' },
  { name: 'points', align: 'center', label: 'Pts', field: 'id' },
];

const filteredFixtures = computed(() => {
  let list = tournamentStore.matches;

  // 1. Filter by Selected Phase Dropdown
  if (selectedPhaseFilter.value !== 'all') {
    list = list.filter((m) => m.phase === selectedPhaseFilter.value);
  }

  // 2. Filter by Text Country search
  if (!filterText.value) return list;
  const filter = filterText.value.toLowerCase();
  return list.filter((m) => {
    const home = m.expand?.home_team?.name.toLowerCase() || '';
    const away = m.expand?.away_team?.name.toLowerCase() || '';
    return home.includes(filter) || away.includes(filter);
  });
});

function translatePhase(phaseKey: string): string {
  switch (phaseKey) {
    case 'group':
      return 'Fase de Grupos';
    case 'r32':
      return '16avos de Final';
    case 'r16':
      return 'Octavos de Final';
    case 'qf':
      return 'Cuartos de Final';
    case 'sf':
      return 'Semifinal';
    case 'final':
      return 'Final';
    default:
      return 'Fase';
  }
}

async function fetchUserGroups() {
  if (!authStore.user?.id) return;
  try {
    const memberships = await pb.collection('group_members_id').getFullList({
      filter: `user = "${authStore.user.id}"`,
      expand: 'prediction_group',
    });

    myGroups.value = memberships
      .map((m) => m.expand?.prediction_group as unknown as PredictionGroup)
      .filter(Boolean);

    if (myGroups.value.length > 0 && myGroups.value[0]) {
      activeGroupId.value = String(myGroups.value[0].id);
    }
  } catch (err: unknown) {
    console.error('Failed to load user groups:', err);
  }
}

async function fetchPredictions() {
  if (!authStore.user?.id || !activeGroupId.value) return;
  loading.value = true;
  try {
    const list = await pb.collection('predictions_id').getFullList({
      filter: `user = "${authStore.user.id}" && prediction_group = "${activeGroupId.value}"`,
    });

    const map: Record<string, Prediction> = {};
    list.forEach((p) => {
      let mId = '';
      if (typeof p.match === 'string') {
        mId = p.match;
      } else if (Array.isArray(p.match) && p.match.length > 0) {
        mId = String(p.match[0]);
      } else if (p.match && typeof p.match === 'object') {
        const matchObj = p.match as { id?: string };
        mId = String(matchObj.id || '');
      }

      if (mId) {
        map[mId] = p as unknown as Prediction;
      }
    });
    predictionsMap.value = map;
  } catch (err: unknown) {
    console.error('Failed loading fixtures predictions:', err);
  } finally {
    loading.value = false;
  }
}

watch(
  [() => authStore.user?.id, () => activeGroupId.value],
  async ([userId, groupId]) => {
    if (userId && groupId) {
      await fetchPredictions();
    }
  },
  { immediate: true },
);

watch(
  () => authStore.user?.id,
  async (newUserId) => {
    if (newUserId) {
      await tournamentStore.fetchMatches();
      await fetchUserGroups();
    }
  },
  { immediate: true },
);
</script>

<style scoped>
.border {
  border: 1px solid rgba(0, 0, 0, 0.1);
}
.border-grey-8 {
  border-color: rgba(255, 255, 255, 0.08) !important;
}
.bracket-col {
  width: 250px;
}
.bracket-match {
  width: 100%;
}
</style>
