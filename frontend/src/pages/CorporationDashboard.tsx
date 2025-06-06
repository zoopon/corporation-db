import React, { useState, useMemo } from 'react';
import { 
  MagnifyingGlassIcon,
  BuildingOfficeIcon,
  MapPinIcon,
  PhoneIcon,
  GlobeAltIcon,
  XMarkIcon,
  CalendarIcon,
  CurrencyYenIcon,
  UsersIcon,
  CreditCardIcon
} from '@heroicons/react/24/outline';
import { useCorporations, useRefreshBases, useCorporation } from '../services/corporationService';
import { useDebounce } from '../hooks/useDebounce';

const CorporationDashboard: React.FC = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedCorporateNumber, setSelectedCorporateNumber] = useState<string>('');
  const [showDetailPanel, setShowDetailPanel] = useState(false);
  
  // 検索キーワードをデバウンス（500ms遅延）
  const debouncedSearchTerm = useDebounce(searchTerm, 500);
  
  // 検索パラメータを作成
  const searchParams = useMemo(() => {
    if (!debouncedSearchTerm.trim()) {
      return undefined;
    }
    
    // 法人番号の形式かどうかをチェック（13桁の数字）
    const isCorporateNumber = /^\d{13}$/.test(debouncedSearchTerm.replace(/\D/g, ''));
    
    if (isCorporateNumber) {
      return { corporateNumber: debouncedSearchTerm.replace(/\D/g, '') };
    } else {
      return { name: debouncedSearchTerm };
    }
  }, [debouncedSearchTerm]);
  
  // 企業一覧取得（検索パラメータ付き）
  const { data: corporations, isLoading: corporationsLoading, error: corporationsError } = useCorporations(searchParams);
  
  // 選択された企業の詳細情報を取得
  const { data: selectedCorporation, isLoading: corporationDetailsLoading } = useCorporation(selectedCorporateNumber);
  
  // 選択された企業の拠点情報（詳細APIから取得）
  const bases = selectedCorporation?.bases || [];
  
  // 拠点情報再取得ミューテーション
  const refreshBasesMutation = useRefreshBases();

  // 企業リスト（APIから直接取得、フィルタリング不要）
  const corporationsList = corporations?.corporations || [];

  const handleCorporationSelect = (corporateNumber: string) => {
    setSelectedCorporateNumber(corporateNumber);
    setShowDetailPanel(true);
  };

  const handleRefreshBases = () => {
    if (selectedCorporateNumber) {
      refreshBasesMutation.mutate(selectedCorporateNumber);
    }
  };

  const closeDetailPanel = () => {
    setShowDetailPanel(false);
    setSelectedCorporateNumber('');
  };

  const formatCurrency = (value: number | undefined | null) => {
    if (!value) return '未設定';
    return new Intl.NumberFormat('ja-JP', {
      style: 'currency',
      currency: 'JPY',
      minimumFractionDigits: 0,
    }).format(value);
  };

  const formatDate = (dateString: string | undefined | null) => {
    if (!dateString) return '未設定';
    return new Date(dateString).toLocaleDateString('ja-JP');
  };

  return (
    <div className="max-w-7xl mx-auto relative">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">企業情報ダッシュボード</h1>
        <p className="text-gray-600">企業情報と拠点データの閲覧・管理</p>
      </div>

      {/* 企業一覧 (フルワイド) */}
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
              {/* 検索中インジケーター */}
              {searchTerm !== debouncedSearchTerm && (
                <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
                </div>
              )}
            </div>
            
            {/* 検索結果の情報 */}
            {searchParams && (
              <div className="mt-2 text-sm text-gray-600">
                {searchParams.name && `企業名: "${searchParams.name}" で検索中`}
                {searchParams.corporateNumber && `法人番号: "${searchParams.corporateNumber}" で検索中`}
              </div>
            )}
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
              <>
                {/* 検索結果件数 */}
                {searchParams && (
                  <div className="px-4 py-2 bg-gray-50 border-b text-sm text-gray-600">
                    検索結果: {corporationsList.length}件
                  </div>
                )}
                
                <div className="divide-y divide-gray-200">
                  {corporationsList.map((corp) => (
                  <div
                    key={corp.id}
                    onClick={() => handleCorporationSelect(corp.corporate_number)}
                    className={`p-4 cursor-pointer hover:bg-blue-50 hover:border-l-2 hover:border-blue-300 transition-all duration-200 ${
                      selectedCorporateNumber === corp.corporate_number 
                        ? 'bg-blue-50 border-l-2 border-blue-500 shadow-sm' 
                        : 'hover:shadow-sm'
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
                
                {corporationsList.length === 0 && !corporationsLoading && (
                  <div className="p-6 text-center text-gray-500">
                    <p>該当する企業が見つかりませんでした</p>
                  </div>
                )}
                </div>
              </>
            )}
          </div>
        </div>

      {/* 右側スライドアウトパネル: 企業詳細 */}
      <div className={`fixed inset-y-0 right-0 z-50 w-full sm:w-96 bg-white shadow-xl transform transition-transform duration-300 ease-in-out ${
        showDetailPanel ? 'translate-x-0' : 'translate-x-full'
      }`}>
        <div className="flex flex-col h-full">
          {/* ヘッダー */}
          <div className="flex items-center justify-between p-6 border-b border-gray-200">
            <h2 className="text-lg font-semibold text-gray-900">企業詳細</h2>
            <button
              onClick={closeDetailPanel}
              className="p-2 rounded-md text-gray-400 hover:text-gray-600 hover:bg-gray-100"
            >
              <XMarkIcon className="h-5 w-5" />
            </button>
          </div>

          {/* コンテンツ */}
          <div className="flex-1 overflow-y-auto">
            {corporationDetailsLoading ? (
              <div className="p-6 text-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                <p className="mt-2 text-gray-500">詳細情報を読み込み中...</p>
              </div>
            ) : selectedCorporation ? (
              <div className="p-6 space-y-6">
                {/* 基本情報 */}
                <div>
                  <h3 className="text-sm font-medium text-gray-900 mb-3">基本情報</h3>
                  <div className="space-y-3">
                    <div>
                      <span className="block text-xs text-gray-500">企業名</span>
                      <span className="text-sm text-gray-900">{selectedCorporation.name}</span>
                    </div>
                    <div>
                      <span className="block text-xs text-gray-500">法人番号</span>
                      <span className="text-sm text-gray-900">{selectedCorporation.corporate_number}</span>
                    </div>
                    {selectedCorporation.kana && (
                      <div>
                        <span className="block text-xs text-gray-500">フリガナ</span>
                        <span className="text-sm text-gray-900">{selectedCorporation.kana}</span>
                      </div>
                    )}
                    {selectedCorporation.name_en && (
                      <div>
                        <span className="block text-xs text-gray-500">英語名</span>
                        <span className="text-sm text-gray-900">{selectedCorporation.name_en}</span>
                      </div>
                    )}
                    {selectedCorporation.location && (
                      <div>
                        <span className="block text-xs text-gray-500">所在地</span>
                        <span className="text-sm text-gray-900">{selectedCorporation.location}</span>
                      </div>
                    )}
                    {selectedCorporation.postal_code && (
                      <div>
                        <span className="block text-xs text-gray-500">郵便番号</span>
                        <span className="text-sm text-gray-900">{selectedCorporation.postal_code}</span>
                      </div>
                    )}
                    <div>
                      <span className="block text-xs text-gray-500">法人状態</span>
                      <span className="text-sm text-gray-900">{selectedCorporation.status}</span>
                    </div>
                  </div>
                </div>

                {/* 代表者情報 */}
                {(selectedCorporation.representative_name || selectedCorporation.representative_position) && (
                  <div>
                    <h3 className="text-sm font-medium text-gray-900 mb-3">代表者情報</h3>
                    <div className="space-y-3">
                      {selectedCorporation.representative_name && (
                        <div>
                          <span className="block text-xs text-gray-500">代表者名</span>
                          <span className="text-sm text-gray-900">{selectedCorporation.representative_name}</span>
                        </div>
                      )}
                      {selectedCorporation.representative_position && (
                        <div>
                          <span className="block text-xs text-gray-500">役職</span>
                          <span className="text-sm text-gray-900">{selectedCorporation.representative_position}</span>
                        </div>
                      )}
                    </div>
                  </div>
                )}

                {/* 企業詳細 */}
                {(selectedCorporation.date_of_establishment || selectedCorporation.founding_year || selectedCorporation.capital_stock || selectedCorporation.employee_number) && (
                  <div>
                    <h3 className="text-sm font-medium text-gray-900 mb-3">企業詳細</h3>
                    <div className="space-y-3">
                      {selectedCorporation.date_of_establishment && (
                        <div className="flex items-center">
                          <CalendarIcon className="h-4 w-4 text-gray-400 mr-2" />
                          <div>
                            <span className="block text-xs text-gray-500">設立年月日</span>
                            <span className="text-sm text-gray-900">{formatDate(selectedCorporation.date_of_establishment)}</span>
                          </div>
                        </div>
                      )}
                      {selectedCorporation.founding_year && (
                        <div className="flex items-center">
                          <CalendarIcon className="h-4 w-4 text-gray-400 mr-2" />
                          <div>
                            <span className="block text-xs text-gray-500">創業年</span>
                            <span className="text-sm text-gray-900">{selectedCorporation.founding_year}年</span>
                          </div>
                        </div>
                      )}
                      {selectedCorporation.capital_stock && (
                        <div className="flex items-center">
                          <CurrencyYenIcon className="h-4 w-4 text-gray-400 mr-2" />
                          <div>
                            <span className="block text-xs text-gray-500">資本金</span>
                            <span className="text-sm text-gray-900">{formatCurrency(selectedCorporation.capital_stock)}</span>
                          </div>
                        </div>
                      )}
                      {selectedCorporation.employee_number && (
                        <div className="flex items-center">
                          <UsersIcon className="h-4 w-4 text-gray-400 mr-2" />
                          <div>
                            <span className="block text-xs text-gray-500">従業員数</span>
                            <span className="text-sm text-gray-900">{selectedCorporation.employee_number}人</span>
                          </div>
                        </div>
                      )}
                      {(selectedCorporation.company_size_male || selectedCorporation.company_size_female) && (
                        <div className="flex items-center">
                          <UsersIcon className="h-4 w-4 text-gray-400 mr-2" />
                          <div>
                            <span className="block text-xs text-gray-500">男女別従業員数</span>
                            <span className="text-sm text-gray-900">
                              男性: {selectedCorporation.company_size_male || 0}人 / 女性: {selectedCorporation.company_size_female || 0}人
                            </span>
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                )}

                {/* 事業情報 */}
                {(selectedCorporation.business_items || selectedCorporation.business_summary) && (
                  <div>
                    <h3 className="text-sm font-medium text-gray-900 mb-3">事業情報</h3>
                    <div className="space-y-3">
                      {selectedCorporation.business_items && (
                        <div>
                          <span className="block text-xs text-gray-500">事業内容</span>
                          <span className="text-sm text-gray-900">{selectedCorporation.business_items}</span>
                        </div>
                      )}
                      {selectedCorporation.business_summary && (
                        <div>
                          <span className="block text-xs text-gray-500">事業概要</span>
                          <span className="text-sm text-gray-900">{selectedCorporation.business_summary}</span>
                        </div>
                      )}
                    </div>
                  </div>
                )}

                {/* 連絡先情報 */}
                {(selectedCorporation.company_url || selectedCorporation.qualification_grade) && (
                  <div>
                    <h3 className="text-sm font-medium text-gray-900 mb-3">その他情報</h3>
                    <div className="space-y-3">
                      {selectedCorporation.company_url && (
                        <div className="flex items-center">
                          <GlobeAltIcon className="h-4 w-4 text-gray-400 mr-2" />
                          <div>
                            <span className="block text-xs text-gray-500">ウェブサイト</span>
                            <a 
                              href={selectedCorporation.company_url} 
                              target="_blank" 
                              rel="noopener noreferrer"
                              className="text-sm text-blue-600 hover:text-blue-800 underline"
                            >
                              {selectedCorporation.company_url}
                            </a>
                          </div>
                        </div>
                      )}
                      {selectedCorporation.qualification_grade && (
                        <div className="flex items-center">
                          <CreditCardIcon className="h-4 w-4 text-gray-400 mr-2" />
                          <div>
                            <span className="block text-xs text-gray-500">資格等級</span>
                            <span className="text-sm text-gray-900">{selectedCorporation.qualification_grade}</span>
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                )}

                {/* メタデータ */}
                <div>
                  <h3 className="text-sm font-medium text-gray-900 mb-3">データ情報</h3>
                  <div className="space-y-3">
                    {selectedCorporation.update_date && (
                      <div>
                        <span className="block text-xs text-gray-500">最終更新日</span>
                        <span className="text-sm text-gray-900">{formatDate(selectedCorporation.update_date)}</span>
                      </div>
                    )}
                    <div>
                      <span className="block text-xs text-gray-500">登録日</span>
                      <span className="text-sm text-gray-900">{formatDate(selectedCorporation.created_at)}</span>
                    </div>
                  </div>
                </div>

                {/* 拠点情報 */}
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-sm font-medium text-gray-900">拠点情報</h3>
                    <button
                      onClick={handleRefreshBases}
                      disabled={refreshBasesMutation.isPending}
                      className="px-3 py-1 text-xs font-medium text-white bg-blue-600 rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {refreshBasesMutation.isPending ? '取得中...' : '拠点情報を取得'}
                    </button>
                  </div>
                  
                  <div className="space-y-3">
                    {refreshBasesMutation.isPending ? (
                      <div className="text-center py-4">
                        <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600 mx-auto"></div>
                        <p className="mt-2 text-xs text-gray-500">拠点情報を取得中...</p>
                      </div>
                    ) : bases && bases.length > 0 ? (
                      bases.map((base) => (
                        <div key={base.id} className="border border-gray-200 rounded-lg p-3">
                          <div className="flex items-start">
                            <MapPinIcon className="h-4 w-4 text-gray-400 mt-1 mr-2 flex-shrink-0" />
                            <div className="flex-1 min-w-0">
                              <h4 className="text-sm font-medium text-gray-900">
                                {base.base_name || '拠点'}
                              </h4>
                              <p className="text-xs text-gray-600 mt-1">
                                {base.location}
                              </p>
                              
                              <div className="mt-2 space-y-1">
                                {base.phone_number && (
                                  <div className="flex items-center text-xs text-gray-500">
                                    <PhoneIcon className="h-3 w-3 mr-1" />
                                    {base.phone_number}
                                  </div>
                                )}
                                {base.data_source_url && (
                                  <div className="flex items-center text-xs text-gray-500">
                                    <GlobeAltIcon className="h-3 w-3 mr-1" />
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
                      ))
                    ) : (
                      <div className="text-center py-4 text-gray-500">
                        <MapPinIcon className="h-8 w-8 mx-auto text-gray-300 mb-2" />
                        <p className="text-xs">拠点情報が見つかりませんでした</p>
                        <p className="text-xs mt-1">「拠点情報を取得」ボタンを押してAIから情報を取得できます</p>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ) : (
              <div className="p-6 text-center text-gray-500">
                <BuildingOfficeIcon className="h-12 w-12 mx-auto text-gray-300 mb-4" />
                <p>企業詳細が見つかりませんでした</p>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* オーバーレイ */}
      {showDetailPanel && (
        <div 
          className="fixed inset-0 bg-black bg-opacity-50 z-40"
          onClick={closeDetailPanel}
        />
      )}
    </div>
  );
};

export default CorporationDashboard;
