import { useEffect, useState } from "react";
import LinearProgress from "@mui/material/LinearProgress";
import { apiGetRequest } from "../../utils/api";

const StarRating = ({ rating }) => {
  const stars = [];
  for (let i = 1; i <= 5; i++) {
    const filled = i <= Math.floor(rating);
    const half = !filled && i === Math.ceil(rating) && rating % 1 !== 0;
    stars.push(
      <span
        key={i}
        className={`text-lg ${filled ? "text-yellow-400" : half ? "text-yellow-300" : "text-gray-300"}`}
      >
        ★
      </span>
    );
  }
  return (
    <span className="flex items-center gap-0.5">
      {stars}
      <span className="ml-1 text-sm text-gray-500">{rating}/5</span>
    </span>
  );
};

const MetricBar = ({ label, value, color }) => (
  <div className="flex flex-col gap-1">
    <div className="flex justify-between text-xs text-gray-500">
      <span>{label}</span>
      <span>{value}%</span>
    </div>
    <LinearProgress
      variant="determinate"
      value={value}
      sx={{
        height: 6,
        borderRadius: 3,
        backgroundColor: "#e5e7eb",
        "& .MuiLinearProgress-bar": { backgroundColor: color, borderRadius: 3 },
      }}
    />
  </div>
);

const ProgressBar = ({ label, current, total, color }) => {
  const pct = total > 0 ? Math.min(Math.round((current / total) * 100), 100) : 0;
  return (
    <div className="flex flex-col gap-1">
      <div className="flex justify-between text-xs text-gray-600">
        <span>{label}</span>
        <span>
          {current} / {total} ({pct}%)
        </span>
      </div>
      <LinearProgress
        variant="determinate"
        value={pct}
        sx={{
          height: 8,
          borderRadius: 4,
          backgroundColor: "#e5e7eb",
          "& .MuiLinearProgress-bar": { backgroundColor: color, borderRadius: 4 },
        }}
      />
    </div>
  );
};

function AppCommunityDashboard() {
  const [details, setDetails] = useState(null);
  const [members, setMembers] = useState([]);
  const [activities, setActivities] = useState([]);
  const [events, setEvents] = useState([]);
  const [targets, setTargets] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchAll = async () => {
      const [detRes, memRes, actRes, evtRes, tgtRes] = await Promise.all([
        apiGetRequest("/community/details"),
        apiGetRequest("/community/members"),
        apiGetRequest("/community/activities"),
        apiGetRequest("/community/events"),
        apiGetRequest("/community/targets"),
      ]);
      if (detRes.success) setDetails(detRes.data);
      if (memRes.success) setMembers(memRes.data);
      if (actRes.success) setActivities(actRes.data);
      if (evtRes.success) setEvents(evtRes.data);
      if (tgtRes.success) setTargets(tgtRes.data);
      setLoading(false);
    };
    fetchAll();
  }, []);

  if (loading) {
    return (
      <div className="w-full p-8 flex justify-center items-center h-64">
        <div className="text-gray-500 animate-pulse">Loading community data…</div>
      </div>
    );
  }

  return (
    <div className="p-5 flex flex-col gap-6 bg-gray-50 min-h-screen">
      {/* Header */}
      <div className="flex items-center gap-3">
        <i className="bx bx-group text-3xl text-primary"></i>
        <div>
          <h1 className="text-2xl font-bold text-gray-800">Community Dashboard</h1>
          <p className="text-sm text-gray-500">Overview of your community's performance and activities</p>
        </div>
      </div>

      {/* Community Overview Card */}
      {details && (
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
          <div className="flex flex-col md:flex-row md:items-center justify-between gap-4">
            <div className="flex items-center gap-4">
              <div className="w-14 h-14 rounded-full bg-primary flex items-center justify-center text-white text-2xl">
                <i className={`bx ${details.icon || "bx-group"}`}></i>
              </div>
              <div>
                <h2 className="text-xl font-bold text-gray-800">{details.name}</h2>
                <p className="text-sm text-gray-500">
                  Established: {details.established_date}
                </p>
                <StarRating rating={details.rating ?? 4.5} />
              </div>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
              <div className="flex flex-col items-center p-3 bg-purple-50 rounded-lg">
                <span className="text-2xl font-bold text-purple-600">{details.total_points}</span>
                <span className="text-xs text-gray-500 mt-1">Total Points</span>
              </div>
              <div className="flex flex-col items-center p-3 bg-blue-50 rounded-lg">
                <span className="text-2xl font-bold text-blue-600">{details.member_count}</span>
                <span className="text-xs text-gray-500 mt-1">Members</span>
              </div>
              <div className="col-span-2 md:col-span-1 flex flex-col gap-2 p-3 bg-gray-50 rounded-lg">
                <MetricBar label="Reliability" value={details.reliability ?? 85} color="#7D53F6" />
                <MetricBar label="Quality" value={details.quality ?? 78} color="#3B82F6" />
                <MetricBar label="Frequency" value={details.frequency ?? 92} color="#10B981" />
              </div>
            </div>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Community Targets */}
        {targets && (
          <div className="lg:col-span-2 bg-white rounded-xl shadow-sm border border-gray-100 p-6 flex flex-col gap-4">
            <h3 className="font-semibold text-gray-700 text-lg border-b pb-2">Community Targets</h3>
            <ProgressBar
              label="Weekly Volume"
              current={targets.weekly_current}
              total={targets.weekly_target}
              color="#7D53F6"
            />
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 mt-2">
              {targets.mandates && targets.mandates.map((m, i) => (
                <div
                  key={i}
                  className={`flex items-center gap-3 p-3 rounded-lg border ${
                    m.completed ? "border-green-200 bg-green-50" : "border-orange-200 bg-orange-50"
                  }`}
                >
                  <i
                    className={`bx ${m.completed ? "bx-check-circle text-green-500" : "bx-time text-orange-400"} text-xl`}
                  ></i>
                  <div>
                    <p className="text-sm font-medium text-gray-700">{m.title}</p>
                    <p className="text-xs text-gray-500">{m.completed ? "Completed" : "Pending"}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Active Members */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 flex flex-col gap-3">
          <h3 className="font-semibold text-gray-700 text-lg border-b pb-2">Active Members</h3>
          {members.length === 0 ? (
            <p className="text-sm text-gray-400 text-center py-4">No active members found</p>
          ) : (
            <div className="flex flex-col gap-3 overflow-auto max-h-64">
              {members.map((m, i) => (
                <div key={i} className="flex items-center gap-3">
                  <div className="w-9 h-9 rounded-full bg-primary flex items-center justify-center text-white text-sm font-semibold flex-shrink-0">
                    {m.name ? m.name.charAt(0).toUpperCase() : "?"}
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-700 truncate">{m.name}</p>
                    <p className="text-xs text-gray-400 truncate">{m.role}</p>
                  </div>
                  <span className="text-xs font-semibold text-purple-600 bg-purple-50 px-2 py-0.5 rounded-full">
                    {m.points} pts
                  </span>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Recent Activities */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 flex flex-col gap-3">
          <h3 className="font-semibold text-gray-700 text-lg border-b pb-2">Recent Activities</h3>
          {activities.length === 0 ? (
            <p className="text-sm text-gray-400 text-center py-4">No recent activities</p>
          ) : (
            <div className="flex flex-col gap-3 overflow-auto max-h-64">
              {activities.map((a, i) => (
                <div key={i} className="flex gap-3 items-start">
                  <div className="w-8 h-8 rounded-full bg-blue-50 flex items-center justify-center flex-shrink-0">
                    <i className={`bx ${a.icon || "bx-bell"} text-blue-500`}></i>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-gray-700">{a.title}</p>
                    <p className="text-xs text-gray-400">{a.description}</p>
                  </div>
                  <span className="text-xs text-gray-400 flex-shrink-0">{a.time}</span>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Events & Announcements */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-100 p-6 flex flex-col gap-3">
          <h3 className="font-semibold text-gray-700 text-lg border-b pb-2">Events &amp; Announcements</h3>
          {events.length === 0 ? (
            <p className="text-sm text-gray-400 text-center py-4">No upcoming events</p>
          ) : (
            <div className="flex flex-col gap-3 overflow-auto max-h-64">
              {events.map((e, i) => (
                <div
                  key={i}
                  className="flex gap-3 items-start p-3 rounded-lg bg-gray-50 border border-gray-100"
                >
                  <div
                    className={`w-10 h-10 rounded-lg flex flex-col items-center justify-center flex-shrink-0 ${
                      e.type === "event" ? "bg-purple-100 text-purple-600" : "bg-amber-100 text-amber-600"
                    }`}
                  >
                    <i className={`bx ${e.type === "event" ? "bx-calendar-event" : "bx-megaphone"} text-lg`}></i>
                  </div>
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-semibold text-gray-700">{e.title}</p>
                    <p className="text-xs text-gray-400">{e.description}</p>
                    <p className="text-xs text-purple-500 mt-1">{e.date}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default AppCommunityDashboard;
