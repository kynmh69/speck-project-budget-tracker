'use client';

import { useCurrentUser } from '@/hooks/use-auth';
import Loading from '@/components/common/loading';

export default function DashboardPage() {
  const { data: user, isLoading } = useCurrentUser();

  if (isLoading) {
    return <Loading />;
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">ダッシュボード</h1>
      
      <div className="bg-white p-6 rounded-lg shadow mb-6">
        <h2 className="text-xl font-semibold mb-4">ようこそ、{user?.name}さん</h2>
        <p className="text-gray-600">
          プロジェクト予算管理システムにログインしています。
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="text-sm text-gray-600 mb-2">総プロジェクト数</div>
          <div className="text-3xl font-bold">0</div>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="text-sm text-gray-600 mb-2">進行中</div>
          <div className="text-3xl font-bold">0</div>
        </div>
        
        <div className="bg-white p-6 rounded-lg shadow">
          <div className="text-sm text-gray-600 mb-2">完了</div>
          <div className="text-3xl font-bold">0</div>
        </div>
      </div>

      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-xl font-semibold mb-4">最近のプロジェクト</h2>
        <p className="text-gray-500 text-center py-8">
          プロジェクトがまだありません
        </p>
      </div>
    </div>
  );
}
