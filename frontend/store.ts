import Vuex, {type StoreOptions} from 'vuex'
import VuexPersist from 'vuex-persist'

// Define the shape of the state
interface RootState {
  mode: string;
  settings: {
    darkTheme: boolean;
  };
}

// Create vuex-persist plugin
const vuexLocalStorage = new VuexPersist<RootState>({
  key: 'vuex', // The key to store the state on in the storage provider.
  storage: window.localStorage, // or window.sessionStorage or localForage
  reducer: (state: RootState) => ({
    settings: state.settings,
  })
})

// Define the store with type annotations
const store: StoreOptions<RootState> = {
  plugins: [vuexLocalStorage.plugin],
  state: {
    mode: "preprocessor",
    settings: {
      darkTheme: false,
    },
  },
  getters: {
    getSettings: (state) => state.settings,
    getThemeVariation: (state) => state.settings.darkTheme ? "dark" : "light",
    getMode: (state) => state.mode,
  },
  mutations: {
    toggleTheme(state, val: boolean) {
      state.settings.darkTheme = val
    },
    setMode(state, val: string) {
      state.mode = val
    },
  },
  actions: {}
}

export default new Vuex.Store<RootState>(store)