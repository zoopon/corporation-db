import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api, type ApiError } from './api';
import { useAppStore, type Corporation } from '../stores/appStore';

// APIレスポンスからアプリケーション型への変換関数
const mapApiCorporationToApp = (apiCorp: any): Corporation => ({
  id: apiCorp.id,
  corporateNumber: apiCorp.corporate_number,
  name: apiCorp.name,
  address: apiCorp.location,
  prefecture: apiCorp.prefecture,
  city: apiCorp.city,
  status: apiCorp.status,
  closeDate: apiCorp.close_date,
  closeCause: apiCorp.close_cause,
  successorCorporateNumber: apiCorp.successor_corporate_number,
  changeDate: apiCorp.change_date,
  assignmentDate: apiCorp.assignment_date,
  latestUpdateDate: apiCorp.latest_update_date,
  enName: apiCorp.name_en,
  enPrefecture: apiCorp.prefecture_en,
  enCity: apiCorp.city_en,
  enAddress: apiCorp.address_en,
  kana: apiCorp.kana,
  bases: apiCorp.bases || [],
});

// 企業一覧取得フック（検索対応版）
export const useCorporations = (searchParams?: { name?: string; corporateNumber?: string }) => {
  const { setCorporations, setCorporationLoading, setCorporationError } = useAppStore();
  
  // クエリキーに検索パラメータを含める
  const queryKey = ['corporations', searchParams?.name, searchParams?.corporateNumber];
  
  return useQuery({
    queryKey,
    queryFn: async () => {
      setCorporationLoading(true);
      
      try {
        // クエリパラメータを構築
        const queryParams = new URLSearchParams();
        if (searchParams?.name && searchParams.name.trim()) {
          queryParams.append('name', searchParams.name.trim());
        }
        if (searchParams?.corporateNumber && searchParams.corporateNumber.trim()) {
          queryParams.append('corporate_number', searchParams.corporateNumber.trim());
        }
        
        const queryString = queryParams.toString();
        const url = queryString ? `/corporations?${queryString}` : '/corporations';
        
        const { data, error } = await api.GET(url as '/corporations');
        
        if (error) {
          const apiError: ApiError = {
            message: error.error || 'Failed to fetch corporations',
            status: 500, // デフォルトステータス
          };
          setCorporationError(apiError.message);
          throw apiError;
        }
        
        if (data) {
          const mappedCorporations = (data.corporations || []).map(mapApiCorporationToApp);
          setCorporations(mappedCorporations);
          setCorporationError(null);
        }
        
        return data;
      } catch (error) {
        const apiError = error as ApiError;
        setCorporationError(apiError.message);
        throw error;
      } finally {
        setCorporationLoading(false);
      }
    },
    staleTime: 5 * 60 * 1000, // 5分間はキャッシュを使用
    // 検索パラメータがある場合はより短いキャッシュ時間
    ...(searchParams?.name || searchParams?.corporateNumber ? { staleTime: 30 * 1000 } : {}),
  });
};

// 特定企業の詳細取得フック
export const useCorporation = (corporateNumber: string) => {
  const { setSelectedCorporation } = useAppStore();
  
  return useQuery({
    queryKey: ['corporation', corporateNumber],
    queryFn: async () => {
      const { data, error } = await api.GET('/corporations/{corporate_number}', {
        params: {
          path: { corporate_number: corporateNumber },
        },
      });
      
      if (error) {
        throw new Error('Failed to fetch corporation details');
      }
      
      if (data) {
        const mappedCorporation = mapApiCorporationToApp(data);
        setSelectedCorporation(mappedCorporation);
      }
      
      return data;
    },
    enabled: !!corporateNumber,
  });
};

// 企業の拠点取得フック
export const useCorporationBases = (corporateNumber: string) => {
  const { setBases, setBasesLoading, setBasesError } = useAppStore();
  
  return useQuery({
    queryKey: ['corporation-bases', corporateNumber],
    queryFn: async () => {
      setBasesLoading(true);
      
      try {
        const { data, error } = await api.POST('/corporations/{corporate_number}/fetch-bases', {
          params: {
            path: { corporate_number: corporateNumber },
          },
        });
        
        if (error) {
          const apiError: ApiError = {
            message: error.error || 'Failed to fetch corporation bases',
            status: 500,
          };
          setBasesError(apiError.message);
          throw apiError;
        }
        
        if (data) {
          // 拠点情報のレスポンスは直接配列ではなく、メタデータを含むオブジェクト
          // 実際の拠点データは関連する企業情報から取得する必要がある
          setBases([]); // 現在は空配列を設定
          setBasesError(null);
        }
        
        return data;
      } catch (error) {
        const apiError = error as ApiError;
        setBasesError(apiError.message);
        throw error;
      } finally {
        setBasesLoading(false);
      }
    },
    enabled: !!corporateNumber,
    staleTime: 10 * 60 * 1000, // 10分間はキャッシュを使用（OpenAI API呼び出しは時間がかかるため）
  });
};

// 拠点情報を再取得するミューテーション
export const useRefreshBases = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (corporateNumber: string) => {
      const { data, error } = await api.POST('/corporations/{corporate_number}/fetch-bases', {
        params: {
          path: { corporate_number: corporateNumber },
        },
      });
      
      if (error) {
        throw new Error('Failed to refresh corporation bases');
      }
      
      return data;
    },
    onSuccess: (_, corporateNumber) => {
      // 関連するクエリを無効化して再取得
      queryClient.invalidateQueries({ queryKey: ['corporation-bases', corporateNumber] });
      queryClient.invalidateQueries({ queryKey: ['corporation', corporateNumber] });
    },
  });
};
