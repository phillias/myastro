import { useState } from 'react';
import { MapContainer, TileLayer } from 'react-leaflet';
import * as Dialog from '@radix-ui/react-dialog';
import * as Tabs from '@radix-ui/react-tabs';
import { Sparkles, X, MapPin, Calendar, Clock, Heart, Sun } from 'lucide-react';
import MapPicker from './MapPicker';

interface Remedy { mantra: string; yantra: string; color: string }
interface Panchangam { tithi: string; nakshatra: string; yoga: string }

interface PlanetPosition {
  name: string; longitude: number; sign: string; degree: number;
  nakshatra: string; pada: number; house: number; is_retrograde: boolean;
}

interface Chart {
  lagna: PlanetPosition; planets: PlanetPosition[];
  house_cusps: number[]; time: string; latitude: number; longitude: number;
}

interface DashaPeriod {
  planet: string; start_date: string; end_date: string; years: number;
}
interface DashaResult {
  birth_time: string; nakshatra: string; nakshatra_index: number;
  pada: number; starting_lord: string; periods: DashaPeriod[];
}

interface JyotishData {
  sign: string; nakshatra: string; nakshatra_pada: number;
  lagana_degree: number; remedy: Remedy; panchanga: Panchangam;
  chart?: Chart; dasha?: DashaResult;
}

interface Aspect {
  planet_a: string; planet_b: string; aspect: string; orb: number;
}
interface VedicDrishti {
  planet: string; from_chart: string; aspects_house: number;
  in_chart: string; significance: string;
}
interface HouseOverlay {
  planet: string; from_chart: string; in_house: number;
  of_chart: string; house_lord: string; meaning: string;
}
interface SynastryResult {
  chart_a: Chart; chart_b: Chart; aspects: Aspect[];
  vedic_drishti: VedicDrishti[]; house_overlays: HouseOverlay[]; summary: string;
}

function PlanetTable({ chart }: { chart: Chart }) {
  const allBodies = [chart.lagna, ...chart.planets];
  return (
    <div className="space-y-1 text-sm">
      <table className="w-full text-xs">
        <thead><tr className="text-slate-400 border-b border-white/10">
          <th className="text-left py-1 pr-1">Body</th>
          <th className="text-left py-1 pr-1">Sign</th>
          <th className="text-right py-1 pr-1">Deg</th>
          <th className="text-right py-1 pr-1">House</th>
          <th className="text-left py-1">Nakshatra</th>
        </tr></thead>
        <tbody>
          {allBodies.map((p) => (
            <tr key={p.name} className="border-b border-white/5">
              <td className="py-1 pr-1 font-medium">
                {p.name}
                {p.is_retrograde && <span className="text-purple-400 ml-0.5">R</span>}
              </td>
              <td className="py-1 pr-1 text-slate-300">{p.sign}</td>
              <td className="py-1 pr-1 text-right text-slate-300">{p.degree.toFixed(1)}°</td>
              <td className="py-1 pr-1 text-right text-slate-300">{p.house}</td>
              <td className="py-1 text-slate-300 truncate max-w-[80px]">{p.nakshatra}</td>
            </tr>
          ))}
        </tbody>
      </table>
      <p className="text-xs text-slate-500 mt-2">R = Retrograde</p>
    </div>
  );
}

function SectionTitle({ icon, text }: { icon: React.ReactNode; text: string }) {
  return (
    <h2 className="text-lg font-semibold mb-4 flex items-center gap-2">
      {icon} {text}
    </h2>
  );
}

export default function App() {
  const [mode, setMode] = useState<'single' | 'compatibility'>('single');

  const [coordsA, setCoordsA] = useState<{ lat: number; lon: number } | null>(null);
  const [dateA, setDateA] = useState('');
  const [timeA, setTimeA] = useState('');

  const [coordsB, setCoordsB] = useState<{ lat: number; lon: number } | null>(null);
  const [dateB, setDateB] = useState('');
  const [timeB, setTimeB] = useState('');

  const [singleResult, setSingleResult] = useState<JyotishData | null>(null);
  const [compatResult, setCompatResult] = useState<SynastryResult | null>(null);
  const [dialogOpen, setDialogOpen] = useState(false);

  // ----- Single chart: fetch /api/calculate -----
  const fetchSingle = async () => {
    if (!coordsA) return;
    setSingleResult(null);
    try {
      const url = `/api/calculate?lat=${coordsA.lat}&lon=${coordsA.lon}&date=${dateA}&time=${timeA}`;
      const res = await fetch(url);
      const data: JyotishData = await res.json();
      setSingleResult(data);
      setDialogOpen(true);
    } catch (err) {
      console.error('Failed to fetch chart:', err);
    }
  };

  // ----- Compatibility: fetch /api/compatibility -----
  const fetchCompatibility = async () => {
    if (!coordsA || !coordsB) return;
    setCompatResult(null);
    try {
      const u = `/api/compatibility?lat1=${coordsA.lat}&lon1=${coordsA.lon}&date1=${dateA}&time1=${timeA}&lat2=${coordsB.lat}&lon2=${coordsB.lon}&date2=${dateB}&time2=${timeB}`;
      const res = await fetch(u);
      const data: SynastryResult = await res.json();
      setCompatResult(data);
      setDialogOpen(true);
    } catch (err) {
      console.error('Failed to fetch compatibility:', err);
    }
  };

  // ----- Map form for one person -----
  const PersonForm = ({
    label, setCoords, date, setDate, time, setTime,
  }: {
    label: string;
    setCoords: (c: { lat: number; lon: number } | null) => void;
    date: string; setDate: (v: string) => void;
    time: string; setTime: (v: string) => void;
  }) => (
    <div className="bg-slate-800 p-6 rounded-2xl border border-white/10 shadow-2xl backdrop-blur-md">
      <h3 className="text-md font-semibold mb-3 flex items-center gap-2 text-purple-300">
        <MapPin size={14} /> {label}
      </h3>
      <div className="space-y-4">
        <div className="h-48 w-full rounded-xl overflow-hidden border border-white/10 relative">
          <MapContainer center={[20, 0]} zoom={2} className="h-full w-full !block">
            <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
            <MapPicker onLocationSelect={(lat, lon) => setCoords({ lat, lon })} />
          </MapContainer>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div className="space-y-1">
            <label className="text-xs text-slate-400 flex items-center gap-1"><Calendar size={12}/> Date</label>
            <input type="date" value={date} onChange={(e) => setDate(e.target.value)}
              className="w-full bg-slate-900 border border-white/10 rounded-lg p-2 text-sm" />
          </div>
          <div className="space-y-1">
            <label className="text-xs text-slate-400 flex items-center gap-1"><Clock size={12}/> Time</label>
            <input type="time" value={time} onChange={(e) => setTime(e.target.value)}
              className="w-full bg-slate-900 border border-white/10 rounded-lg p-2 text-sm" />
          </div>
        </div>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-slate-900 text-white p-6">
      <header className="flex items-center gap-3 mb-8">
        <Sparkles className="text-purple-400" />
        <h1 className="text-3xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-purple-400 to-indigo-400">
          Jyotish PWA
        </h1>
      </header>

      {/* ---- Mode toggle ---- */}
      <div className="max-w-md mx-auto mb-6">
        <div className="flex bg-slate-800 rounded-xl p-1 border border-white/10">
          <button
            onClick={() => setMode('single')}
            className={`flex-1 flex items-center justify-center gap-2 py-2 px-4 rounded-lg text-sm font-medium transition ${
              mode === 'single' ? 'bg-purple-600 text-white shadow-lg' : 'text-slate-400 hover:text-white'
            }`}
          >
            <Sun size={16} /> Single Chart
          </button>
          <button
            onClick={() => setMode('compatibility')}
            className={`flex-1 flex items-center justify-center gap-2 py-2 px-4 rounded-lg text-sm font-medium transition ${
              mode === 'compatibility' ? 'bg-purple-600 text-white shadow-lg' : 'text-slate-400 hover:text-white'
            }`}
          >
            <Heart size={16} /> Compatibility
          </button>
        </div>
      </div>

      <main className="max-w-md mx-auto space-y-6">
        {mode === 'single' ? (
          <>
            <div className="bg-slate-800 p-6 rounded-2xl border border-white/10 shadow-2xl backdrop-blur-md">
              <SectionTitle icon={<MapPin size={18} />} text="Enter Birth Details" />
              <div className="space-y-4">
                <div className="h-48 w-full rounded-xl overflow-hidden border border-white/10 relative">
                  <MapContainer center={[20, 0]} zoom={2} className="h-full w-full !block">
                    <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
                    <MapPicker onLocationSelect={(lat, lon) => setCoordsA({ lat, lon })} />
                  </MapContainer>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div className="space-y-1">
                    <label className="text-xs text-slate-400 flex items-center gap-1"><Calendar size={12}/> Date</label>
                    <input type="date" value={dateA} onChange={(e) => setDateA(e.target.value)}
                      className="w-full bg-slate-900 border border-white/10 rounded-lg p-2 text-sm" />
                  </div>
                  <div className="space-y-1">
                    <label className="text-xs text-slate-400 flex items-center gap-1"><Clock size={12}/> Time</label>
                    <input type="time" value={timeA} onChange={(e) => setTimeA(e.target.value)}
                      className="w-full bg-slate-900 border border-white/10 rounded-lg p-2 text-sm" />
                  </div>
                </div>
                <button
                  onClick={fetchSingle}
                  disabled={!coordsA}
                  className="w-full bg-gradient-to-r from-purple-600 to-indigo-600 p-3 rounded-xl font-bold hover:from-purple-500 hover:to-indigo-500 transition shadow-lg disabled:opacity-50"
                >
                  <Sun size={16} className="inline mr-2" />Get Jyotish Remedies
                </button>
              </div>
            </div>
          </>
        ) : (
          <>
            <PersonForm
              label="Person A"
              setCoords={setCoordsA}
              date={dateA} setDate={setDateA} time={timeA} setTime={setTimeA}
            />
            <PersonForm
              label="Person B"
              setCoords={setCoordsB}
              date={dateB} setDate={setDateB} time={timeB} setTime={setTimeB}
            />
            <button
              onClick={fetchCompatibility}
              disabled={!coordsA || !coordsB}
              className="w-full bg-gradient-to-r from-purple-600 to-indigo-600 p-3 rounded-xl font-bold hover:from-purple-500 hover:to-indigo-500 transition shadow-lg disabled:opacity-50 flex items-center justify-center gap-2"
            >
              <Heart size={18} /> Get Compatibility
            </button>
          </>
        )}
      </main>

      {/* ---- Single chart result dialog ---- */}
      <Dialog.Root open={dialogOpen && mode === 'single' && !!singleResult} onOpenChange={setDialogOpen}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50" />
          <Dialog.Content className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-slate-800 border border-white/10 p-6 rounded-2xl shadow-2xl w-[380px] max-h-[80vh] overflow-y-auto z-[60]">
            <Dialog.Title className="text-xl font-bold mb-4 flex items-center gap-2">
              <Sun size={18} className="text-purple-400" /> Your Cosmic Report
            </Dialog.Title>
            <Dialog.Description className="hidden">Detailed Jyotish insights</Dialog.Description>
            {singleResult && (
              <Tabs.Root defaultValue="remedy">
                <Tabs.List className="flex border-b border-white/10 mb-4 overflow-x-auto">
                  <Tabs.Trigger value="remedy" className="px-3 py-2 text-xs data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Natal</Tabs.Trigger>
                  <Tabs.Trigger value="panchanga" className="px-3 py-2 text-xs data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Panchanga</Tabs.Trigger>
                  <Tabs.Trigger value="planets" className="px-3 py-2 text-xs data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Planets</Tabs.Trigger>
                  <Tabs.Trigger value="dasha" className="px-3 py-2 text-xs data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Dashas</Tabs.Trigger>
                </Tabs.List>

                <Tabs.Content value="remedy" className="space-y-2 text-sm">
                  <p className="text-slate-400">Sign: <span className="text-white">{singleResult.sign}</span></p>
                  <div className="mt-4 p-4 bg-white/5 rounded-xl border border-white/5">
                    <p className="font-semibold text-purple-300">Mantra</p>
                    <p className="italic mb-2">{singleResult.remedy.mantra}</p>
                    <p className="font-semibold text-purple-300">Color</p>
                    <p>{singleResult.remedy.color}</p>
                  </div>
                </Tabs.Content>

                <Tabs.Content value="panchanga" className="space-y-2 text-sm">
                  <p className="text-slate-400">Tithi: <span className="text-white">{singleResult.panchanga.tithi}</span></p>
                  <p className="text-slate-400">Nakshatra: <span className="text-white">{singleResult.panchanga.nakshatra}</span></p>
                </Tabs.Content>

                <Tabs.Content value="planets" className="text-sm">
                  {singleResult.chart ? (
                    <PlanetTable chart={singleResult.chart} />
                  ) : (
                    <p className="text-slate-400 italic">Planet data not available</p>
                  )}
                </Tabs.Content>

                <Tabs.Content value="dasha" className="text-sm max-h-[350px] overflow-y-auto">
                  {singleResult.dasha ? (
                    <div className="space-y-2">
                      <p className="text-xs text-slate-400">
                        Birth Nakshatra: <span className="text-white">{singleResult.dasha.nakshatra}</span>
                        , Pada: <span className="text-white">{singleResult.dasha.pada + 1}</span>
                        , Starting Lord: <span className="text-purple-300">{singleResult.dasha.starting_lord}</span>
                      </p>
                      <table className="w-full text-xs">
                        <thead><tr className="text-slate-400 border-b border-white/10">
                          <th className="text-left py-1 pr-1">Period</th>
                          <th className="text-left py-1 pr-1">From</th>
                          <th className="text-left py-1 pr-1">To</th>
                          <th className="text-right py-1">Yrs</th>
                        </tr></thead>
                        <tbody>
                          {singleResult.dasha.periods.map((p, i) => (
                            <tr key={i} className="border-b border-white/5">
                              <td className="py-1 pr-1 font-medium">{p.planet}</td>
                              <td className="py-1 pr-1 text-slate-300">{p.start_date}</td>
                              <td className="py-1 pr-1 text-slate-300">{p.end_date}</td>
                              <td className="py-1 text-right text-slate-300">{p.years}</td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  ) : (
                    <p className="text-slate-400 italic">Dasha data not available</p>
                  )}
                </Tabs.Content>
              </Tabs.Root>
            )}
            <Dialog.Close className="absolute top-4 right-4 text-slate-400 hover:text-white"><X size={20}/></Dialog.Close>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>

      {/* ---- Compatibility result dialog ---- */}
      <Dialog.Root open={dialogOpen && mode === 'compatibility' && !!compatResult} onOpenChange={setDialogOpen}>
        <Dialog.Portal>
          <Dialog.Overlay className="fixed inset-0 bg-black/70 backdrop-blur-sm z-50" />
          <Dialog.Content className="fixed top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-slate-800 border border-white/10 p-6 rounded-2xl shadow-2xl w-[420px] max-h-[85vh] overflow-y-auto z-[60]">
            <Dialog.Title className="text-xl font-bold mb-4 flex items-center gap-2">
              <Heart size={18} className="text-purple-400" /> Compatibility Report
            </Dialog.Title>
            <Dialog.Description className="hidden">Synastry compatibility analysis</Dialog.Description>
            {compatResult && (
              <Tabs.Root defaultValue="overview">
                <Tabs.List className="flex border-b border-white/10 mb-4 overflow-x-auto">
                  <Tabs.Trigger value="overview" className="px-3 py-2 text-xs data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Overview</Tabs.Trigger>
                  <Tabs.Trigger value="aspects" className="px-3 py-2 text-xs data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Aspects</Tabs.Trigger>
                  <Tabs.Trigger value="drishti" className="px-3 py-2 text-xs data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Drishti</Tabs.Trigger>
                  <Tabs.Trigger value="overlays" className="px-3 py-2 text-xs data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Overlays</Tabs.Trigger>
                  <Tabs.Trigger value="charts" className="px-3 py-2 text-xs data-[state=active]:border-b-2 data-[state=active]:border-purple-400">Charts</Tabs.Trigger>
                </Tabs.List>

                {/* Overview */}
                <Tabs.Content value="overview" className="space-y-2 text-sm">
                  <div className="p-3 bg-white/5 rounded-xl">
                    <p className="text-xs text-purple-400 mb-1">Chart A — Lagna in {compatResult.chart_a.lagna.sign} ({compatResult.chart_a.lagna.degree.toFixed(1)}°)</p>
                    <p className="text-xs text-purple-400 mb-1">Chart B — Lagna in {compatResult.chart_b.lagna.sign} ({compatResult.chart_b.lagna.degree.toFixed(1)}°)</p>
                  </div>
                  <div className="p-3 bg-white/5 rounded-xl">
                    <p className="text-xs text-slate-400 mb-1">Aspect count</p>
                    <div className="grid grid-cols-5 gap-1 text-center text-xs">
                      {['conjunction','sextile','square','trine','opposition'].map((a) => {
                        const n = compatResult.aspects.filter((x) => x.aspect === a).length;
                        return <div key={a} className="bg-slate-900 rounded p-1"><span className="block text-purple-400 font-bold">{n}</span>{a.slice(0,4)}</div>;
                      })}
                    </div>
                  </div>
                  <pre className="whitespace-pre-wrap text-xs text-slate-300 leading-relaxed bg-slate-900 p-3 rounded-xl max-h-[300px] overflow-y-auto">
                    {compatResult.summary}
                  </pre>
                </Tabs.Content>

                {/* Aspects */}
                <Tabs.Content value="aspects" className="space-y-1 text-sm max-h-[400px] overflow-y-auto">
                  {compatResult.aspects.length === 0 ? (
                    <p className="text-slate-400 italic">No Ptolemaic aspects detected.</p>
                  ) : (
                    <table className="w-full text-xs">
                      <thead><tr className="text-slate-400 border-b border-white/10">
                        <th className="text-left py-1 pr-1">A</th>
                        <th className="text-left py-1 pr-1">Aspect</th>
                        <th className="text-left py-1 pr-1">B</th>
                        <th className="text-right py-1">Orb</th>
                      </tr></thead>
                      <tbody>
                        {compatResult.aspects.map((asp, i) => (
                          <tr key={i} className="border-b border-white/5">
                            <td className="py-1 pr-1 font-medium">{asp.planet_a}</td>
                            <td className="py-1 pr-1">
                              <span className={`px-1.5 py-0.5 rounded text-xs ${
                                ['trine','sextile'].includes(asp.aspect)
                                  ? 'bg-green-900/40 text-green-300'
                                  : asp.aspect === 'conjunction'
                                  ? 'bg-blue-900/40 text-blue-300'
                                  : 'bg-red-900/40 text-red-300'
                              }`}>{asp.aspect}</span>
                            </td>
                            <td className="py-1 pr-1 font-medium">{asp.planet_b}</td>
                            <td className="py-1 text-right text-slate-400">{asp.orb.toFixed(1)}°</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  )}
                </Tabs.Content>

                {/* Vedic Drishti */}
                <Tabs.Content value="drishti" className="text-sm max-h-[400px] overflow-y-auto">
                  {compatResult.vedic_drishti.length === 0 ? (
                    <p className="text-slate-400 italic">No Vedic graha drishti cross-chart.</p>
                  ) : (
                    <div className="space-y-2">
                      {compatResult.vedic_drishti.map((d, i) => (
                        <div key={i} className="p-2 bg-white/5 rounded-lg text-xs">
                          <span className="text-purple-300 font-medium">{d.planet}</span>
                          <span className="text-slate-400"> (chart {d.from_chart}) → </span>
                          <span className="text-slate-300">House {d.aspects_house}</span>
                          <span className="text-slate-400"> in chart {d.in_chart}</span>
                          <p className="text-slate-500 mt-0.5">{d.significance}</p>
                        </div>
                      ))}
                    </div>
                  )}
                </Tabs.Content>

                {/* House Overlays */}
                <Tabs.Content value="overlays" className="text-sm max-h-[400px] overflow-y-auto">
                  {compatResult.house_overlays.length === 0 ? (
                    <p className="text-slate-400 italic">No house overlays.</p>
                  ) : (
                    <table className="w-full text-xs">
                      <thead><tr className="text-slate-400 border-b border-white/10">
                        <th className="text-left py-1 pr-1">Planet</th>
                        <th className="text-right py-1 pr-1">House</th>
                        <th className="text-left py-1 pr-1">In</th>
                        <th className="text-left py-1">Lord</th>
                      </tr></thead>
                      <tbody>
                        {compatResult.house_overlays.map((o, i) => (
                          <tr key={i} className="border-b border-white/5">
                            <td className="py-1 pr-1 font-medium">{o.planet}</td>
                            <td className="py-1 pr-1 text-right text-slate-300">{o.in_house}</td>
                            <td className="py-1 pr-1 text-slate-400">Chart {o.of_chart}</td>
                            <td className="py-1 text-slate-300">{o.house_lord}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  )}
                  {compatResult.house_overlays.length > 0 && (
                    <div className="mt-3 space-y-1">
                      <p className="text-xs text-slate-500 font-semibold mb-1">Notable overlays:</p>
                      {compatResult.house_overlays.filter(o => [1,5,7,10].includes(o.in_house)).map((o, i) => (
                        <p key={i} className="text-xs text-slate-400">{o.meaning}</p>
                      ))}
                    </div>
                  )}
                </Tabs.Content>

                {/* Charts (Planet Tables) */}
                <Tabs.Content value="charts" className="space-y-4 text-sm max-h-[400px] overflow-y-auto">
                  <div>
                    <p className="text-xs text-purple-400 font-semibold mb-2">Chart A — Planets</p>
                    <PlanetTable chart={compatResult.chart_a} />
                  </div>
                  <div>
                    <p className="text-xs text-purple-400 font-semibold mb-2">Chart B — Planets</p>
                    <PlanetTable chart={compatResult.chart_b} />
                  </div>
                </Tabs.Content>
              </Tabs.Root>
            )}
            <Dialog.Close className="absolute top-4 right-4 text-slate-400 hover:text-white"><X size={20}/></Dialog.Close>
          </Dialog.Content>
        </Dialog.Portal>
      </Dialog.Root>
    </div>
  );
}
