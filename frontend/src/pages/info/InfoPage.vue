<template>
  <q-page class="q-pa-md">
    <div class="row justify-center">
      <div class="col-12 col-md-10 col-lg-8">
        <q-card flat bordered :class="$q.dark.isActive ? 'bg-grey-10' : 'bg-white'">
          <q-card-section>
            <div class="text-h4 text-weight-bold text-primary q-mb-md">
              Mundial 2026 — Plataforma de Predicciones (Prode)
            </div>
            <div class="text-body1 q-mb-lg">
              Bienvenido a la plataforma de predicciones de fútbol <strong>Mundial 2026</strong>, diseñada para la Copa Mundial de la FIFA 2026. Esta plataforma permite a los usuarios unirse a grupos, predecir resultados de partidos y escalar en las tablas de posiciones de los grupos en tiempo real.
            </div>

            <q-separator class="q-my-lg" />

            <div class="text-h5 text-weight-bold text-primary q-mb-md">
              ⚽ Reglas y Sistema de Asignación de Puntos
            </div>
            <div class="text-body1 q-mb-md">
              Los puntos se otorgan dinámicamente según los resultados de los partidos y la precisión. Los partidos de las etapas finales tienen mayor peso para recompensar el éxito en las predicciones del final del torneo.
            </div>

            <div class="text-h6 text-weight-bold q-mt-lg q-mb-sm">
              1. Lógica de Puntos en Fase de Grupos
            </div>
            <ul class="q-mb-md">
              <li><strong>Solo Resultado Correcto</strong> (Victoria, Empate o Derrota): <strong>+2 Puntos</strong></li>
              <li><strong>Bonus por Resultado Exacto</strong> (ej: predecir 2-1 y el partido termina 2-1): <strong>+1 Punto</strong> (Total <strong>3 Puntos</strong> otorgados)</li>
            </ul>

            <div class="text-h6 text-weight-bold q-mt-lg q-mb-sm">
              2. Multiplicadores de Puntos en Fase Eliminatoria
            </div>
            <div class="text-body2 q-mb-sm">
              A medida que aumenta la tensión, los partidos otorgan ganancias de puntos progresivamente mayores:
            </div>
            <q-table
              :rows="pointsTable"
              :columns="pointsColumns"
              flat
              bordered
              dense
              hide-bottom
              :class="$q.dark.isActive ? 'bg-grey-9' : 'bg-grey-1'"
            />

            <q-separator class="q-my-lg" />

            <div class="text-h5 text-weight-bold text-primary q-mb-md">
              📝 Cómo Participar
            </div>

            <div class="text-h6 text-weight-bold q-mt-md q-mb-sm">
              Paso 1: Conexión de Cuenta
            </div>
            <ul class="q-mb-md">
              <li>Abre la aplicación, registra un perfil o conéctate instantáneamente usando <strong>Autenticación de Google</strong>.</li>
            </ul>

            <div class="text-h6 text-weight-bold q-mt-md q-mb-sm">
              Paso 2: Unirse o Crear un Grupo de Predicciones
            </div>
            <ul class="q-mb-md">
              <li><strong>Crear un Grupo</strong>: Desde la pestaña <strong>Grupos</strong>, haz clic en <em>Crear Grupo</em>. Recibirás instantáneamente un token de invitación de 8 caracteres.</li>
              <li><strong>Unirse a un Grupo</strong>: Pega un token de invitación de 8 caracteres de amigos para unirte a su tabla de posiciones privada.</li>
            </ul>

            <div class="text-h6 text-weight-bold q-mt-md q-mb-sm">
              Paso 3: Ingresar Predicciones (Panel Mobile-First)
            </div>
            <ul class="q-mb-md">
              <li>Desliza a través de los Grupos A a L usando el control deslizante superior.</li>
              <li>Escribe tus predicciones. Los teclados numéricos nativos para móviles se activan automáticamente.</li>
              <li>Los números negativos y caracteres no numéricos están bloqueados.</li>
              <li>Haz clic en <strong>Guardar Predicciones del Grupo</strong> para confirmar todas las predicciones de ese grupo en una sola transacción. ¡La página avanzará automáticamente al siguiente grupo!</li>
            </ul>

            <div class="text-h6 text-weight-bold q-mt-md q-mb-sm">
              Paso 4: Seguir las Llaves y Clasificaciones
            </div>
            <ul class="q-mb-md">
              <li>Navega a la pestaña <strong>Llaves</strong> para ver tu progreso global o inspeccionar las etapas eliminatorias vacías.</li>
              <li>Alterna el <strong>icono Sol/Luna</strong> en el encabezado para cambiar entre el tema celeste de la Camiseta Local de Argentina y el tema carbón eléctrico de la Camiseta Alternativa de Argentina.</li>
            </ul>
          </q-card-section>
        </q-card>
      </div>
    </div>
  </q-page>
</template>

<script setup lang="ts">
const pointsTable = [
  { phase: 'Fase de Grupos', correctResult: '+2 Pts', exactBonus: '+1 Pt', maxPoints: '3 Pts' },
  { phase: 'Octavos / Dieciseisavos', correctResult: '+3 Pts', exactBonus: '+2 Pts', maxPoints: '5 Pts' },
  { phase: 'Cuartos / Semifinales', correctResult: '+5 Pts', exactBonus: '+3 Pts', maxPoints: '8 Pts' },
  { phase: 'Final / Tercer Puesto', correctResult: '+8 Pts', exactBonus: '+5 Pts', maxPoints: '13 Pts' },
];

const pointsColumns = [
  { name: 'phase', label: 'Fase del Torneo', field: 'phase', align: 'left' as const },
  { name: 'correctResult', label: 'Resultado Correcto (V/E/D)', field: 'correctResult', align: 'center' as const },
  { name: 'exactBonus', label: 'Bonus Resultado Exacto', field: 'exactBonus', align: 'center' as const },
  { name: 'maxPoints', label: 'Máximo de Puntos Posibles', field: 'maxPoints', align: 'center' as const },
];
</script>
