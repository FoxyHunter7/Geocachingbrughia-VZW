<?php

namespace Database\Seeders;

use App\Models\Geocache;
use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;

class GeocacheSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/geocaches.csv'), ';');

        $model = new Geocache();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
