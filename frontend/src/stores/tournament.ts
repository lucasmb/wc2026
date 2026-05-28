import { defineStore } from 'pinia';
import { ref } from 'vue';
import { pb } from '@/boot/pocketbase';
import type { Match, Team } from 'src/types';

export const useTournamentStore = defineStore('tournament', () => {
  const matches = ref<Match[]>([]);
  const teams = ref<Team[]>([]);
  const currentPhase = ref('group');

  async function fetchMatches() {
    const rawList = await pb.collection('matches_id').getFullList({
      expand: 'home_team,away_team,winner',
      sort: 'kickoff',
    });
    matches.value = rawList as unknown as Match[];
  }

  async function fetchTeams() {
    const rawTeams = await pb.collection('teams_id').getFullList({
      sort: 'group_code,group_rank',
    });
    teams.value = rawTeams as unknown as Team[];
  }

  // SSE real-time sync wrapper
  function subscribeToMatches() {
    // Explicitly handle the subscribe floating promise
    void pb.collection('matches_id').subscribe('*', (e) => {
      const updatedMatch = e.record as unknown as Match;
      const index = matches.value.findIndex((m) => m.id === updatedMatch.id);
      if (index !== -1) {
        const oldMatch = matches.value[index];
        if (oldMatch && oldMatch.expand) {
          // Construct object without any potentially undefined fields
          const cleanExpand: Match['expand'] = {};
          if (oldMatch.expand.home_team) {
            cleanExpand.home_team = oldMatch.expand.home_team;
          }
          if (oldMatch.expand.away_team) {
            cleanExpand.away_team = oldMatch.expand.away_team;
          }
          if (oldMatch.expand.winner) {
            cleanExpand.winner = oldMatch.expand.winner;
          }
          updatedMatch.expand = cleanExpand;
        }
        matches.value[index] = updatedMatch;
      }
    });
  }

  function unsubscribeFromMatches() {
    // Explicitly handle unsubscribe promise
    void pb.collection('matches_id').unsubscribe();
  }
  return {
    matches,
    teams,
    currentPhase,
    fetchMatches,
    fetchTeams,
    subscribeToMatches,
    unsubscribeFromMatches,
  };
});
