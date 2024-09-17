<?php

namespace Database\Seeders;

use App\Models\Geocache;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class GeocacheSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('geocaches')->truncate();

        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/geocaches.csv'), ';');

        $model = new Geocache();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
