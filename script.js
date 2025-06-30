import http from 'k6/http';
import { check, sleep, group } from 'k6';

// --- 1. テスト全体のオプション設定 ---
export const options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '1m', target: 10 },
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    'http_req_failed': ['rate<0.01'],
    'checks': ['rate>0.98'], // 98%以上のチェックが成功すること
  },
};

// --- 2. 定数と事前準備 ---
const BASE_URL = 'http://localhost:8080';

// ★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★
// ★ ここに、事前にデータベースに登録済みのテストユーザー情報を入力してください ★
const TEST_USER_EMAIL = '23610119yt@stu.yamato-u.ac.jp'; 
const TEST_USER_PASSWORD = '0513Yuuki';           
// ★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★★


// --- 3. setup関数: テスト実行前に一度だけ実行 ---
export function setup() {
  console.log('テスト準備: 認証トークンを取得します...');
  const loginRes = http.post(`${BASE_URL}/login`, JSON.stringify({
    email: TEST_USER_EMAIL,
    password: TEST_USER_PASSWORD,
  }), { headers: { 'Content-Type': 'application/json' } });

  if (loginRes.status !== 200 || !loginRes.json('token')) {
    throw new Error('ログインに失敗しました。テストユーザー情報を確認してください。');
  }
  const authToken = loginRes.json('token');
  console.log('ログイン成功。テストを開始します。');
  
  return { token: authToken };
}


// --- 4. default関数: 各仮想ユーザーが繰り返し実行するメイン処理 ---
export default function (data) {
  const token = data.token;
  const authHeaders = {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  };

  const uniqueId = `${__VU}-${__ITER}`; // 毎回ユニークなデータを作るためのID

  // === シナリオ1: 会社リスト (CompanyList) のCRUD ===
  group('CompanyList CRUD', function () {
    let companyId;

    // 1. Create - 必須項目(Company, Member)をすべて含める
    const createCompanyPayload = JSON.stringify({
      Company: `Test Company ${uniqueId}`, // not null
      Occupation: "Software Engineer",
      Member: 100, // not null
      Selection: "1st Interview",
      Intern: true,
    });
    const createRes = http.post(`${BASE_URL}/company_lists`, createCompanyPayload, authHeaders);
    check(createRes, { '[Company] Create success (201)': (r) => r.status === 201 });
    if (createRes.status === 201 && createRes.json('ID')) {
      companyId = createRes.json('ID');
    }

    sleep(1);

    if (companyId) { // Createが成功した場合のみ後続を実行
      // 2. Update
      const updateCompanyPayload = JSON.stringify({ Company: `Updated Company ${uniqueId}` });
      http.put(`${BASE_URL}/company_lists/${companyId}`, updateCompanyPayload, authHeaders);

      sleep(1);

      // 3. Delete
      http.del(`${BASE_URL}/company_lists/${companyId}`, null, authHeaders);
    }
  });

  sleep(2);

  // === シナリオ2: インターンシップ (Internship) のCRUD ===
  group('Internship CRUD', function () {
    // 必須項目がないモデルなので、ペイロードは任意
    let internshipId;
    const createInternPayload = JSON.stringify({
      Title: `Awesome Internship ${uniqueId}`,
      Company: `Intern Inc.`,
      Dailystart: 9,
      Dailyfinish: 18,
      Content: "Develop a new feature.",
      Selection: "Pending",
      Joined: false,
    });
    const createRes = http.post(`${BASE_URL}/internships`, createInternPayload, authHeaders);
    check(createRes, { '[Internship] Create success (201)': (r) => r.status === 201 });
    if (createRes.status === 201 && createRes.json('ID')) {
        internshipId = createRes.json('ID');
    }

    if(internshipId) {
        // ... Update, Deleteも同様に追加 ...
    }
  });

  sleep(2);

  // === シナリオ3: 掲示板 (Post, Comment, Like) のCRUD ===
  group('Post & Comment & Like', function () {
    let postId;

    // 1. Create Post - 必須項目(Title, Content, DisplayName)をすべて含める
    const createPostPayload = JSON.stringify({
      title: `Test Post ${uniqueId}`, // not null
      content: "This is the content of the test post.", // not null
      display_name: `User${__VU}`, // not null
    });
    const createRes = http.post(`${BASE_URL}/posts`, createPostPayload, authHeaders);
    check(createRes, { '[Post] Create success (201)': (r) => r.status === 201 });
    if (createRes.status === 201 && createRes.json('ID')) {
      postId = createRes.json('ID');
    }

    sleep(1);

    if (postId) { // Post作成が成功した場合のみ後続を実行
      // 2. Create Comment - 必須項目(Content, DisplayName)をすべて含める
      const commentPayload = JSON.stringify({
        content: "This is a great comment!", // not null
        display_name: `User${__VU}`, // not null
      });
      const commentRes = http.post(`${BASE_URL}/posts/${postId}/comments`, commentPayload, authHeaders);
      check(commentRes, { '[Comment] Create success (201)': (r) => r.status === 201 });
      
      sleep(1);

      // 3. Like & Unlike
      const likeRes = http.post(`${BASE_URL}/posts/${postId}/like`, null, authHeaders);
      check(likeRes, { '[Like] Success (200)': (r) => r.status === 200 });
      sleep(0.5);
      const unlikeRes = http.del(`${BASE_URL}/posts/${postId}/like`, null, authHeaders);
      check(unlikeRes, { '[Unlike] Success (200)': (r) => r.status === 200 });
      
      sleep(1);

      // 4. Delete Post
      http.del(`${BASE_URL}/posts/${postId}`, null, authHeaders);
    }
  });

  sleep(3);
}