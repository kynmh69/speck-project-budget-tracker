export default function HomePage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <div className="z-10 max-w-5xl w-full items-center justify-between font-mono text-sm">
        <h1 className="text-4xl font-bold text-center mb-8">
          プロジェクト予算管理システム
        </h1>
        <p className="text-center text-lg mb-4">
          ソフトウェア開発プロジェクトの工数予実管理と収支管理
        </p>
        <div className="mt-8 grid grid-cols-1 md:grid-cols-3 gap-4">
          <div className="p-6 border rounded-lg">
            <h2 className="text-xl font-semibold mb-2">工数予実管理</h2>
            <p className="text-sm">タスクごとの予定工数と実績工数を記録し、差異を可視化</p>
          </div>
          <div className="p-6 border rounded-lg">
            <h2 className="text-xl font-semibold mb-2">収支管理</h2>
            <p className="text-sm">プロジェクトの売上、コスト、利益を自動計算</p>
          </div>
          <div className="p-6 border rounded-lg">
            <h2 className="text-xl font-semibold mb-2">グラフ可視化</h2>
            <p className="text-sm">予実比較、収支推移をグラフで表示</p>
          </div>
        </div>
        <div className="mt-8 text-center text-sm text-gray-500">
          <p>Status: Setup Phase - Coming Soon</p>
        </div>
      </div>
    </main>
  );
}
