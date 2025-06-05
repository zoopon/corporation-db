import React, { useState } from 'react';
import { 
  MagnifyingGlassIcon,
  BuildingOfficeIcon,
  MapPinIcon,
  PhoneIcon,
  GlobeAltIcon
} from '@heroicons/react/24/outline';
import { useCorporations, useRefreshBases } from '../services/corporationService';

const CorporationDashboard: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCorporateNumber, setSelectedCorporateNumber] = useState<string>('');
  
  // 企業一覧取得
  const { data: corporations, isLoading: corporationsLoading, error: corporationsError } = useCorporations();
  
  // 選択された企業の詳細情報を取得（拠点情報も含む）
  const selectedCorp = corporations?.corporations?.find(corp => corp.corporate_number === selectedCorporateNumber);
  const bases = selectedCorp?.bases || [];
  
  // 拠点情報再取得ミューテーション
  const refreshBasesMutation = useRefreshBases();

  // 検索でフィルタリングされた企業リスト
  const filteredCorporations = corporations?.corporations?.filter(corp =>
    corp.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    corp.corporate_number.includes(searchTerm)
  ) || [];

  const handleCorporationSelect = (corporateNumber: string) => {
    setSelectedCorporateNumber(corporateNumber);
  };

  const handleRefreshBases = () => {
    if (selectedCorporateNumber) {
      refreshBasesMutation.mutate(selectedCorporateNumber);
    }
  };

  return (
    <div className="max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">企業情報ダッシュボード</h1>
        <p className="text-gray-600">企業情報と拠点データの閲覧・管理</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* 左パネル: 企業一覧 */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="p-6 border-b border-gray-200">
            <h2 className="text-xl font-semibold text-gray-900 mb-4">企業一覧</h2>
            
            {/* 検索バー */}
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <MagnifyingGlassIcon className="h-5 w-5 text-gray-400" />
              </div>
              <input
                type="text"
                placeholder="企業名または法人番号で検索..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
          </div>

          <div className="max-h-96 overflow-y-auto">
            {corporationsLoading ? (
              <div className="p-6 text-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                <p className="mt-2 text-gray-500">読み込み中...</p>
              </div>
            ) : corporationsError ? (
              <div className="p-6 text-center text-red-600">
                <p>エラーが発生しました: {corporationsError.message}</p>
              </div>
            ) : (
              <div className="divide-y divide-gray-200">
                {filteredCorporations.map((corp) => (
                  <div
                    key={corp.id}
                    onClick={() => handleCorporationSelect(corp.corporate_number)}
                    className={`p-4 cursor-pointer hover:bg-gray-50 transition-colors ${
                      selectedCorporateNumber === corp.corporate_number ? 'bg-blue-50 border-r-2 border-blue-500' : ''
                    }`}
                  >
                    <div className="flex items-start">
                      <BuildingOfficeIcon className="h-5 w-5 text-gray-400 mt-1 mr-3 flex-shrink-0" />
                      <div className="flex-1 min-w-0">
                        <h3 className="text-sm font-medium text-gray-900 truncate">
                          {corp.name}
                        </h3>
                        <p className="text-xs text-gray-500 mt-1">
                          法人番号: {corp.corporate_number}
                        </p>
                        {corp.location && (
                          <p className="text-xs text-gray-500 mt-1 truncate">
                            {corp.location}
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                ))}
                
                {filteredCorporations.length === 0 && !corporationsLoading && (
                  <div className="p-6 text-center text-gray-500">
                    <p>該当する企業が見つかりませんでした</p>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>

        {/* 右パネル: 拠点情報 */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200">
          <div className="p-6 border-b border-gray-200">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-semibold text-gray-900">拠点情報</h2>
              {selectedCorporateNumber && (
                <button
                  onClick={handleRefreshBases}
                  disabled={refreshBasesMutation.isPending || corporationsLoading}
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {refreshBasesMutation.isPending || corporationsLoading ? '取得中...' : '拠点情報を取得'}
                </button>
              )}
            </div>
          </div>

          <div className="max-h-96 overflow-y-auto">
            {!selectedCorporateNumber ? (
              <div className="p-6 text-center text-gray-500">
                <BuildingOfficeIcon className="h-12 w-12 mx-auto text-gray-300 mb-4" />
                <p>企業を選択してください</p>
              </div>
            ) : refreshBasesMutation.isPending ? (
              <div className="p-6 text-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                <p className="mt-2 text-gray-500">拠点情報を取得中...</p>
              </div>
            ) : bases && bases.length > 0 ? (
              <div className="divide-y divide-gray-200">
                {bases.map((base) => (
                  <div key={base.id} className="p-4">
                    <div className="flex items-start">
                      <MapPinIcon className="h-5 w-5 text-gray-400 mt-1 mr-3 flex-shrink-0" />
                      <div className="flex-1 min-w-0">
                        <h3 className="text-sm font-medium text-gray-900">
                          {base.base_name || '拠点'}
                        </h3>
                        <p className="text-sm text-gray-600 mt-1">
                          {base.location}
                        </p>
                        
                        <div className="mt-2 space-y-1">
                          {base.phone_number && (
                            <div className="flex items-center text-xs text-gray-500">
                              <PhoneIcon className="h-4 w-4 mr-2" />
                              {base.phone_number}
                            </div>
                          )}
                          {base.data_source_url && (
                            <div className="flex items-center text-xs text-gray-500">
                              <GlobeAltIcon className="h-4 w-4 mr-2" />
                              <a 
                                href={base.data_source_url} 
                                target="_blank" 
                                rel="noopener noreferrer"
                                className="text-blue-600 hover:text-blue-800 truncate"
                              >
                                情報源
                              </a>
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="p-6 text-center text-gray-500">
                <MapPinIcon className="h-12 w-12 mx-auto text-gray-300 mb-4" />
                <p>拠点情報が見つかりませんでした</p>
                <p className="text-sm mt-2">「拠点情報を取得」ボタンを押して、OpenAI APIから情報を取得できます</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default CorporationDashboard;
