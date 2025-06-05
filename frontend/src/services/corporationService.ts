import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api, type ApiError } from './api';
import { useAppStore } from '../stores/appStore';

// 企業一覧取得フック
export const useCorporations = () => {
  const { setCorporations, setCorporationLoading, setCorporationError } = useAppStore();
  
  return useQuery({
    queryKey: ['corporations'],
    queryFn: async () => {
      setCorporationLoading(true);
      
      try {
        const { data, error } = await api.GET('/corporations');
        
        if (error) {
          const apiError: ApiError = {
            message: 'Failed to fetch corporations',
            status: error.status,
          };
          setCorporationError(apiError.message);
          throw apiError;
        }
        
        if (data) {
          setCorporations(data);
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
        setSelectedCorporation(data);
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
            message: 'Failed to fetch corporation bases',
            status: error.status,
          };
          setBasesError(apiError.message);
          throw apiError;
        }
        
        if (data) {
          setBases(data);
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
    onSuccess: (data, corporateNumber) => {
      // 関連するクエリを無効化して再取得
      queryClient.invalidateQueries({ queryKey: ['corporation-bases', corporateNumber] });
      queryClient.invalidateQueries({ queryKey: ['corporation', corporateNumber] });
    },
  });
};
