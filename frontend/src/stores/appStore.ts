import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

// 企業データの型（OpenAPI生成型から取得）
export interface Corporation {
  id: string;
  corporateNumber: string;
  name: string;
  address?: string;
  prefecture?: string;
  city?: string;
  status?: string;
  closeDate?: string;
  closeCause?: string;
  successorCorporateNumber?: string;
  changeDate?: string;
  assignmentDate?: string;
  latestUpdateDate?: string;
  enName?: string;
  enPrefecture?: string;
  enCity?: string;
  enAddress?: string;
  kana?: string;
  bases?: Base[];
}

// 拠点データの型
export interface Base {
  id: string;
  corporationId: string;
  name: string;
  address: string;
  postalCode?: string;
  prefectureCode?: string;
  cityCode?: string;
  streetNumber?: string;
  buildingName?: string;
  phoneNumber?: string;
  faxNumber?: string;
  sourceURL?: string;
}

// アプリケーションの状態管理
interface AppState {
  // 企業関連
  corporations: Corporation[];
  selectedCorporation: Corporation | null;
  corporationLoading: boolean;
  corporationError: string | null;
  
  // 拠点関連
  bases: Base[];
  basesLoading: boolean;
  basesError: string | null;
  
  // UI状態
  sidebarOpen: boolean;
  
  // アクション
  setCorporations: (corporations: Corporation[]) => void;
  setSelectedCorporation: (corporation: Corporation | null) => void;
  setCorporationLoading: (loading: boolean) => void;
  setCorporationError: (error: string | null) => void;
  
  setBases: (bases: Base[]) => void;
  setBasesLoading: (loading: boolean) => void;
  setBasesError: (error: string | null) => void;
  
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
}

export const useAppStore = create<AppState>()(
  devtools(
    (set) => ({
      // 初期状態
      corporations: [],
      selectedCorporation: null,
      corporationLoading: false,
      corporationError: null,
      
      bases: [],
      basesLoading: false,
      basesError: null,
      
      sidebarOpen: true,
      
      // 企業関連アクション
      setCorporations: (corporations) => set({ corporations }),
      setSelectedCorporation: (corporation) => set({ selectedCorporation: corporation }),
      setCorporationLoading: (loading) => set({ corporationLoading: loading }),
      setCorporationError: (error) => set({ corporationError: error }),
      
      // 拠点関連アクション
      setBases: (bases) => set({ bases }),
      setBasesLoading: (loading) => set({ basesLoading: loading }),
      setBasesError: (error) => set({ basesError: error }),
      
      // UI関連アクション
      toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
      setSidebarOpen: (open) => set({ sidebarOpen: open }),
    }),
    {
      name: 'app-storage', // localStorage key
    }
  )
);
